package users

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	SecretKey []byte
	db        *gorm.DB
	logger    *zap.Logger
	validate  *validator.Validate
}

type UserServiceParams struct {
	Logger   *zap.Logger
	Validate *validator.Validate
	DB       *gorm.DB
}

func NewUserService(p UserServiceParams) *UserService {

	p.DB.AutoMigrate(&UserModel{})

	return &UserService{
		logger:   p.Logger,
		validate: p.Validate,
		db:       p.DB,
	}
}

func (u *UserService) FindUserByEmail(email string) (*UserModel, error) {
	var user UserModel
	result := u.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *UserService) CreateUser(email string, password string) (*UserModel, error) {

	hash, err := hashPassword(password)

	fmt.Println("hash", hash)

	if err != nil {
		return nil, err
	}

	user := &UserModel{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hash),
	}

	if err := u.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
