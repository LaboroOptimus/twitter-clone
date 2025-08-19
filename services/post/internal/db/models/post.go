package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoincrement"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description" gorm:"not null;type:varchar(300)"`
	ReplyToPostID *uint          `json:"replyToPostId,omitempty" gorm:"index"`
	AuthorID      uint           `json:"authorId" gorm:"not null;index"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"not null"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Image struct {
	ID     uint   `gorm:"primaryKey"`
	PostID uint   `gorm:"not null;index"`
	URL    string `gorm:"not null"`
	Post   Post   `gorm:"constraint:OnDelete:CASCADE;"`
}
