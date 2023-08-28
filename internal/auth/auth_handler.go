package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	users "github.com/tomdoestech/goth/internal/user"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *AuthService
	userService *users.UserService
	validate    *validator.Validate
	logger      *zap.Logger
}

type AuthHandlerParams struct {
	AuthService *AuthService
	UserService *users.UserService
	Validate    *validator.Validate
	Logger      *zap.Logger
}

type loginData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type registrationData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

func NewAuthHandler(p AuthHandlerParams) *AuthHandler {
	return &AuthHandler{
		authService: p.AuthService,
		userService: p.UserService,
		validate:    p.Validate,
		logger:      p.Logger,
	}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")
	data := loginData{
		Email:    email,
		Password: password,
	}

	err := a.validate.Struct(&data)
	if err != nil {
		// Handle validation errors
		handleValidationErrors(w, err)
		return
	}

	user, err := a.userService.FindUserByEmail(data.Email)

	if err != nil {
		a.logger.Error("Error finding user", zap.Error(err))
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Authenticate user
	err = a.authService.VerifyPassword(user.Password, data.Password)

	if err != nil {
		a.logger.Error("Error authenticating user", zap.Error(err))
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := a.authService.GenerateToken(user)

	if err != nil {
		a.logger.Error("Error generating token", zap.Error(err))
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// set the cookie
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "token", Value: token, Expires: expiration, Path: "/"}

	http.SetCookie(w, &cookie)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func handleValidationErrors(w http.ResponseWriter, err error) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Handle other error types
		http.Error(w, "Validation error", http.StatusBadRequest)
		return
	}

	// Convert validationErrors to an array of error messages
	errorMessages := make([]string, 0, len(validationErrors))
	for _, err := range validationErrors {
		errorMessages = append(errorMessages, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
	}

	// Return the array of error messages with status 401
	w.Header().Set("Content-Type", "text/html charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	errorMessagesHTML := ""
	for _, errorMessage := range errorMessages {
		errorMessagesHTML += fmt.Sprintf("<li>%s</li>", errorMessage)
	}

	fmt.Fprintf(w, "<h1>Validation error</h1><ul>%s</ul>", errorMessagesHTML)

}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")
	data := registrationData{
		Email:    email,
		Password: password,
	}

	err := a.validate.Struct(&data)
	if err != nil {

		fmt.Println(err)

		// Handle validation errors
		handleValidationErrors(w, err)
		return
	}

	_, err = a.userService.CreateUser(data.Email, data.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// return html
	w.Header().Set("Content-Type", "text/html charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	fmt.Fprintf(w, "<h1>Registration successful</h1><p>Go to <a href=\"/login\">login</a></p>")
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	// set the cookie
	expiration := time.Now().Add(-365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "token", Value: "", Expires: expiration, Path: "/"}

	http.SetCookie(w, &cookie)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
