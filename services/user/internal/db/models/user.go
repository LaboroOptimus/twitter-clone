package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uint `gorm:"primaryKey"`
	Password          string
	Email             string `gorm:"uniqueIndex:uni_users_email"`
	Nickname          string `gorm:"uniqueIndex:uni_users_nickname"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Bio               string `gorm:"type:text"`
	AvatarURL         string `gorm:"type:text"`
	SubscribersAmount uint   `gorm:"not null;default:0;check:subscribers_amount >= 0"`
}

func NewUser(nickname, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Nickname: nickname,
		Email:    email,
		Password: string(hashedPassword),
	}, nil
}
