package userprofile_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bobfive1/user-management-api/internal/app"
	"github.com/bobfive1/user-management-api/internal/config"
	"github.com/bobfive1/user-management-api/internal/logger"
	"github.com/bobfive1/user-management-api/internal/userprofile"
	"github.com/bobfive1/user-management-api/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserProfileHandler_Create(t *testing.T) {

	gin.SetMode(gin.TestMode)

	type testcase struct {
		name           string
		request        []byte
		repoReturn     userprofile.UserProfile
		errorReturn    error
		wantStatusCode int
	}

	testcases := []testcase{
		{
			name: "UserProfileHandler_Create_Success",
			request: []byte(`
			{
				"user_id":"bobbydev04",
				"password":"psdsystem",
				"first_name":"tawat",
				"last_name":"test",
				"address":"123/123",
				"birthdate":"2005-01-01",
				"email":"bb@gmail.com"
			}
			`),
			repoReturn: userprofile.UserProfile{
				UserID:    "bob01",
				Password:  "password",
				FirstName: "hanlee",
				LastName:  "kiki",
			},
			errorReturn:    nil,
			wantStatusCode: http.StatusOK,
		},
	}

	config, _ := config.LoadConfig()
	logger.ApplyConfig(config.App.Name, config.Logging.Level)
	validation.Init()

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			repo := userprofile.NewUserProfileMockMock()
			repo.On("Create").Return(tc.repoReturn, tc.errorReturn)

			service := userprofile.NewUserProfileService(repo)
			serv := app.ApiServer(config, service).Serv()

			body := bytes.NewBuffer(tc.request)
			req, _ := http.NewRequest("POST", "/api/v1/userprofiles", body)

			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			serv.Handler.ServeHTTP(w, req)
			assert.Equal(t, tc.wantStatusCode, w.Result().StatusCode)
		})
	}

}
