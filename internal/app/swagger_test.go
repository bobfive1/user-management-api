package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bobfive1/user-management-api/internal/config"
	"github.com/bobfive1/user-management-api/internal/logger"
	"github.com/bobfive1/user-management-api/internal/userprofile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userProfileRepoNoop struct{}

func (userProfileRepoNoop) Create(context.Context, userprofile.UserProfile) (userprofile.UserProfileDisplay, error) {
	return userprofile.UserProfileDisplay{}, nil
}

func (userProfileRepoNoop) List(context.Context) ([]userprofile.UserProfileDisplay, error) {
	return nil, nil
}

func (userProfileRepoNoop) GetByID(context.Context, string) (userprofile.UserProfileDisplay, error) {
	return userprofile.UserProfileDisplay{}, nil
}

func (userProfileRepoNoop) GetByIDWithPassword(context.Context, string) (userprofile.UserProfile, error) {
	return userprofile.UserProfile{}, nil
}

func (userProfileRepoNoop) Update(context.Context, string, userprofile.UserProfile) (userprofile.UserProfileDisplay, error) {
	return userprofile.UserProfileDisplay{}, nil
}

func (userProfileRepoNoop) Delete(context.Context, string) error {
	return nil
}

func TestSwaggerRoutes(t *testing.T) {
	require.NoError(t, logger.ApplyConfig("user-management-api-test", "debug"))

	api := ApiServer(&config.AppConfig{
		ServerAPI: config.ServerAPIConfig{
			Port:         ":0",
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
			IdleTimeout:  time.Second,
		},
	}, userprofile.NewUserProfileService(userProfileRepoNoop{}))

	t.Run("openapi json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.json", nil)
		w := httptest.NewRecorder()

		api.Serv().Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

		var spec map[string]any
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &spec))
		assert.Equal(t, "3.0.3", spec["openapi"])
	})

	t.Run("swagger ui", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger", nil)
		w := httptest.NewRecorder()

		api.Serv().Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		assert.Contains(t, w.Body.String(), "/swagger/openapi.json")
	})
}
