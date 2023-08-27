package users

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model

	ID        uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Email    string `gorm:"uniqueIndex" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"-"`
}

func (UserModel) TableName() string {
	return "users"
}
