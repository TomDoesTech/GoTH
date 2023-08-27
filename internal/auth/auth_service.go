package auth

import (
	"errors"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	users "github.com/tomdoestech/goth/internal/user"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	SecretKey []byte
	logger    *zap.Logger
	validate  *validator.Validate
	tokenAuth *jwtauth.JWTAuth
}

type AuthServiceParams struct {
	Logger    *zap.Logger
	SecretKey []byte
	Validate  *validator.Validate
	TokenAuth *jwtauth.JWTAuth
}

func NewAuthService(p AuthServiceParams) *AuthService {
	return &AuthService{
		SecretKey: p.SecretKey,
		logger:    p.Logger,
		validate:  p.Validate,
		tokenAuth: p.TokenAuth,
	}
}

// VerifyPassword checks if a provided password matches the hashed password
func (a *AuthService) VerifyPassword(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func (as *AuthService) FindUserByEmail(email string) (*users.UserModel, error) {
	user := &users.UserModel{
		Email: email,
	}

	return user, nil
}

func (a *AuthService) GenerateToken(user *users.UserModel) (string, error) {

	payload := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"exp": time.Now().Add(time.Hour *
			24).Unix(), // Token expires in 24 hours
	}

	_, tokenString, err := a.tokenAuth.Encode(payload)

	// tokenString, err := token.SignedString(a.SecretKey)

	if err != nil {
		a.logger.Error("Error generating token", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}

func (as *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return as.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
