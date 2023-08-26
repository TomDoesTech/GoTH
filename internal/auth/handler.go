package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func NewAuthHandler(p AuthHandlerParams) *AuthHandler {
	return &AuthHandler{
		authService: p.AuthService,
		userService: p.UserService,
		validate:    p.Validate,
		logger:      p.Logger,
	}
}

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data loginData

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
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

	// Return token to the client
	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type registrationData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(errorMessages)
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data registrationData

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.validate.Struct(&data)
	if err != nil {
		// Handle validation errors
		handleValidationErrors(w, err)
		return
	}

	user, err := a.userService.CreateUser(data.Email, data.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(interface{}(
		map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	))
}
