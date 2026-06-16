package userprofile_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	errInt "github.com/bobfive1/user-management-api/internal/error"
	"github.com/bobfive1/user-management-api/internal/logger"
	"github.com/bobfive1/user-management-api/internal/userprofile"
	"github.com/bobfive1/user-management-api/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userProfileRepoStub struct {
	t *testing.T

	create              func(context.Context, userprofile.UserProfile) (userprofile.UserProfileDisplay, error)
	list                func(context.Context) ([]userprofile.UserProfileDisplay, error)
	getByID             func(context.Context, string) (userprofile.UserProfileDisplay, error)
	getByIDWithPassword func(context.Context, string) (userprofile.UserProfile, error)
	update              func(context.Context, string, userprofile.UserProfile) (userprofile.UserProfileDisplay, error)
	delete              func(context.Context, string) error
}

func (s *userProfileRepoStub) Create(ctx context.Context, request userprofile.UserProfile) (userprofile.UserProfileDisplay, error) {
	if s.create == nil {
		s.t.Fatal("unexpected Create call")
	}
	return s.create(ctx, request)
}

func (s *userProfileRepoStub) List(ctx context.Context) ([]userprofile.UserProfileDisplay, error) {
	if s.list == nil {
		s.t.Fatal("unexpected List call")
	}
	return s.list(ctx)
}

func (s *userProfileRepoStub) GetByID(ctx context.Context, userID string) (userprofile.UserProfileDisplay, error) {
	if s.getByID == nil {
		s.t.Fatal("unexpected GetByID call")
	}
	return s.getByID(ctx, userID)
}

func (s *userProfileRepoStub) GetByIDWithPassword(ctx context.Context, userID string) (userprofile.UserProfile, error) {
	if s.getByIDWithPassword == nil {
		s.t.Fatal("unexpected GetByIDWithPassword call")
	}
	return s.getByIDWithPassword(ctx, userID)
}

func (s *userProfileRepoStub) Update(ctx context.Context, userID string, request userprofile.UserProfile) (userprofile.UserProfileDisplay, error) {
	if s.update == nil {
		s.t.Fatal("unexpected Update call")
	}
	return s.update(ctx, userID, request)
}

func (s *userProfileRepoStub) Delete(ctx context.Context, userID string) error {
	if s.delete == nil {
		s.t.Fatal("unexpected Delete call")
	}
	return s.delete(ctx, userID)
}

func newUserProfileTestRouter(t *testing.T, repo userprofile.UserProfileRepository) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	require.NoError(t, logger.ApplyConfig("user-management-api-test", "debug"))
	validation.Init()

	router := gin.New()
	router.Use(errInt.MiddlewareErrorHandler())

	service := userprofile.NewUserProfileService(repo)
	handler := userprofile.NewUserProfileHandler(service)
	handler.RegisterRoutes(router.Group("/api/v1"))

	return router
}

func performRequest(t *testing.T, router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	req, err := http.NewRequest(method, path, bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	return body
}

func validCreateBody() string {
	return `{
		"user_id": "bobbydev04",
		"password": "psdsystem",
		"first_name": "tawat",
		"last_name": "test",
		"address": "123/123",
		"birthdate": "2005-01-01",
		"email": "bb@gmail.com"
	}`
}

func validUpdateBody() string {
	return `{
		"password": "newpassword",
		"first_name": "updated",
		"last_name": "profile",
		"address": "321/321",
		"birthdate": "2005-01-01",
		"email": "updated@gmail.com"
	}`
}

func TestUserProfileHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			create: func(_ context.Context, request userprofile.UserProfile) (userprofile.UserProfileDisplay, error) {
				assert.Equal(t, "bobbydev04", request.UserID)
				assert.Equal(t, "tawat", request.FirstName)
				assert.NotEqual(t, "psdsystem", request.Password)
				assert.True(t, userprofile.CheckPasswordHash("psdsystem", request.Password))
				return userprofile.UserProfileDisplay{
					UserID:    request.UserID,
					FirstName: request.FirstName,
					LastName:  request.LastName,
					Address:   request.Address,
					Birthdate: request.Birthdate,
					Email:     request.Email,
				}, nil
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPost, "/api/v1/userprofiles", validCreateBody())

		assert.Equal(t, http.StatusOK, w.Code)
		body := decodeBody(t, w)
		assert.Equal(t, "200", body["error_code"])
		assert.Equal(t, "Success", body["error_message"])
		data := body["data"].(map[string]any)
		assert.Equal(t, "bobbydev04", data["user_id"])
		assert.NotContains(t, data, "password")
	})

	t.Run("validation error", func(t *testing.T) {
		repo := &userProfileRepoStub{t: t}
		body := `{"user_id":"bobbydev04","first_name":"tawat","last_name":"test"}`

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPost, "/api/v1/userprofiles", body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "400", response["error_code"])
		assert.Equal(t, "Field validation Error", response["error_message"])
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			create: func(context.Context, userprofile.UserProfile) (userprofile.UserProfileDisplay, error) {
				return userprofile.UserProfileDisplay{}, errors.New("insert failed")
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPost, "/api/v1/userprofiles", validCreateBody())

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "500", response["error_code"])
	})
}

func TestUserProfileHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			list: func(context.Context) ([]userprofile.UserProfileDisplay, error) {
				return []userprofile.UserProfileDisplay{
					{UserID: "bob01", FirstName: "Bob", LastName: "One"},
					{UserID: "bob02", FirstName: "Bob", LastName: "Two"},
				}, nil
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodGet, "/api/v1/userprofiles", "")

		assert.Equal(t, http.StatusOK, w.Code)
		body := decodeBody(t, w)
		data := body["data"].([]any)
		assert.Len(t, data, 2)
		assert.Equal(t, "bob01", data[0].(map[string]any)["user_id"])
		assert.NotContains(t, data[0].(map[string]any), "password")
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			list: func(context.Context) ([]userprofile.UserProfileDisplay, error) {
				return nil, errors.New("select failed")
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodGet, "/api/v1/userprofiles", "")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "500", response["error_code"])
	})
}

func TestUserProfileHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			getByID: func(_ context.Context, userID string) (userprofile.UserProfileDisplay, error) {
				assert.Equal(t, "bob01", userID)
				return userprofile.UserProfileDisplay{UserID: userID, FirstName: "Bob", LastName: "One"}, nil
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodGet, "/api/v1/userprofiles/bob01", "")

		assert.Equal(t, http.StatusOK, w.Code)
		body := decodeBody(t, w)
		data := body["data"].(map[string]any)
		assert.Equal(t, "bob01", data["user_id"])
	})

	t.Run("not found", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			getByID: func(context.Context, string) (userprofile.UserProfileDisplay, error) {
				return userprofile.UserProfileDisplay{}, pgx.ErrNoRows
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodGet, "/api/v1/userprofiles/missing", "")

		assert.Equal(t, http.StatusNotFound, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "404", response["error_code"])
		assert.Equal(t, userprofile.UserProfileNotExist, response["error_message"])
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			getByID: func(context.Context, string) (userprofile.UserProfileDisplay, error) {
				return userprofile.UserProfileDisplay{}, errors.New("select failed")
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodGet, "/api/v1/userprofiles/bob01", "")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "500", response["error_code"])
	})
}

func TestUserProfileHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			update: func(_ context.Context, userID string, request userprofile.UserProfile) (userprofile.UserProfileDisplay, error) {
				assert.Equal(t, "bob01", userID)
				assert.Equal(t, "updated", request.FirstName)
				assert.True(t, userprofile.CheckPasswordHash("newpassword", request.Password))
				return userprofile.UserProfileDisplay{
					UserID:    userID,
					FirstName: request.FirstName,
					LastName:  request.LastName,
					Address:   request.Address,
					Birthdate: request.Birthdate,
					Email:     request.Email,
				}, nil
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPut, "/api/v1/userprofiles/bob01", validUpdateBody())

		assert.Equal(t, http.StatusOK, w.Code)
		body := decodeBody(t, w)
		data := body["data"].(map[string]any)
		assert.Equal(t, "bob01", data["user_id"])
		assert.Equal(t, "updated", data["first_name"])
	})

	t.Run("validation error", func(t *testing.T) {
		repo := &userProfileRepoStub{t: t}
		body := `{"password":"newpassword","first_name":"updated"}`

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPut, "/api/v1/userprofiles/bob01", body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "400", response["error_code"])
	})

	t.Run("not found", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			update: func(context.Context, string, userprofile.UserProfile) (userprofile.UserProfileDisplay, error) {
				return userprofile.UserProfileDisplay{}, pgx.ErrNoRows
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPut, "/api/v1/userprofiles/missing", validUpdateBody())

		assert.Equal(t, http.StatusNotFound, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, userprofile.UserProfileNotExist, response["error_message"])
	})
}

func TestUserProfileHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			delete: func(_ context.Context, userID string) error {
				assert.Equal(t, "bob01", userID)
				return nil
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodDelete, "/api/v1/userprofiles/bob01", "")

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			delete: func(context.Context, string) error {
				return pgx.ErrNoRows
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodDelete, "/api/v1/userprofiles/missing", "")

		assert.Equal(t, http.StatusNotFound, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, userprofile.UserProfileNotExist, response["error_message"])
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			delete: func(context.Context, string) error {
				return errors.New("delete failed")
			},
		}

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodDelete, "/api/v1/userprofiles/bob01", "")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "500", response["error_code"])
	})
}

func TestUserProfileHandler_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hashedPassword := userprofile.HashPassword("psdsystem")
		repo := &userProfileRepoStub{
			t: t,
			getByIDWithPassword: func(_ context.Context, userID string) (userprofile.UserProfile, error) {
				assert.Equal(t, "bob01", userID)
				return userprofile.UserProfile{
					UserID:    userID,
					Password:  hashedPassword,
					FirstName: "Bob",
					LastName:  "One",
				}, nil
			},
		}
		body := `{"user_id":"bob01","password":"psdsystem"}`

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPost, "/api/v1/userprofiles/login", body)

		assert.Equal(t, http.StatusOK, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "200", response["error_code"])
		data := response["data"].(map[string]any)
		assert.Equal(t, "bob01", data["user_id"])
		assert.NotContains(t, data, "password")
	})

	t.Run("wrong password", func(t *testing.T) {
		hashedPassword := userprofile.HashPassword("psdsystem")
		repo := &userProfileRepoStub{
			t: t,
			getByIDWithPassword: func(context.Context, string) (userprofile.UserProfile, error) {
				return userprofile.UserProfile{UserID: "bob01", Password: hashedPassword}, nil
			},
		}
		body := `{"user_id":"bob01","password":"wrongpassword"}`

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPost, "/api/v1/userprofiles/login", body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "400", response["error_code"])
		assert.Equal(t, userprofile.LoginFailMessage, response["error_message"])
	})

	t.Run("user not found", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			getByIDWithPassword: func(context.Context, string) (userprofile.UserProfile, error) {
				return userprofile.UserProfile{}, pgx.ErrNoRows
			},
		}
		body := `{"user_id":"missing","password":"psdsystem"}`

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPost, "/api/v1/userprofiles/login", body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, userprofile.LoginFailMessage, response["error_message"])
	})

	t.Run("validation error", func(t *testing.T) {
		repo := &userProfileRepoStub{t: t}
		body := `{"user_id":"bob01"}`

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPost, "/api/v1/userprofiles/login", body)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "Field validation Error", response["error_message"])
	})

	t.Run("repository error", func(t *testing.T) {
		repo := &userProfileRepoStub{
			t: t,
			getByIDWithPassword: func(context.Context, string) (userprofile.UserProfile, error) {
				return userprofile.UserProfile{}, errors.New("select failed")
			},
		}
		body := `{"user_id":"bob01","password":"psdsystem"}`

		w := performRequest(t, newUserProfileTestRouter(t, repo), http.MethodPost, "/api/v1/userprofiles/login", body)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		response := decodeBody(t, w)
		assert.Equal(t, "500", response["error_code"])
	})
}
