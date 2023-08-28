//go:build unit
// +build unit

package auth

import (
	"io"
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	users "github.com/tomdoestech/goth/internal/user"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TempFilename(t testing.TB) string {
	f, err := os.CreateTemp("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name() + ".db"
}

func TestRegister(t *testing.T) {

	testCases := []struct {
		description          string
		formData             url.Values
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			description: "create user",
			formData: url.Values{
				"email":    {"test@example.com"},
				"password": {"password"},
			},
			expectedStatusCode:   201,
			expectedResponseBody: "<h1>Registration successful</h1><p>Go to <a href=\"/login\">login</a></p>",
		},
		{
			description: "invalid email",
			formData: url.Values{
				"email":    {"test@example"},
				"password": {"password"},
			},
			expectedStatusCode:   400,
			expectedResponseBody: "<h1>Validation error</h1><ul><li>Email is email</li></ul>",
		},
		{
			description: "invalid password",
			formData: url.Values{
				"email":    {"test@example.com"},
				"password": {"1"},
			},
			expectedStatusCode:   400,
			expectedResponseBody: "<h1>Validation error</h1><ul><li>Password is min</li></ul>",
		},
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	filename := TempFilename(t)

	defer os.Remove(filename)

	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var tokenAuth *jwtauth.JWTAuth

	// use a single instance of Validate, it caches struct info
	var validate *validator.Validate

	for _, tc := range testCases {
		validate = validator.New()

		usersService := users.NewUserService(users.UserServiceParams{
			Logger:   logger,
			Validate: validate,
			DB:       db,
		})
		authService := NewAuthService(AuthServiceParams{
			Logger:    logger,
			SecretKey: []byte("secret"),
			TokenAuth: tokenAuth,
		})

		authHandler := NewAuthHandler(
			AuthHandlerParams{
				AuthService: authService,
				UserService: usersService,
				Validate:    validate,
				Logger:      logger,
			},
		)

		t.Run(tc.description, func(t *testing.T) {

			assert := assert.New(t)

			encodedFormData := tc.formData.Encode()
			req := httptest.NewRequest("POST", "/register", strings.NewReader(encodedFormData))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			authHandler.Register(w, req)

			assert.Equal(tc.expectedStatusCode, w.Code)

			res := w.Result()

			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}

			assert.Equal(tc.expectedResponseBody, string(data))
		})
	}

}
