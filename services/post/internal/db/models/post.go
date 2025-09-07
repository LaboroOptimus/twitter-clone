package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID            uint           `gorm:"index:idx_posts_author_created_id,priority:3"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description" gorm:"not null;type:varchar(300)"`
	ReplyToPostID *uint          `json:"replyToPostId,omitempty" gorm:"index"`
	AuthorID      uint           `gorm:"index:idx_posts_author_created_id,priority:1"`
	CreatedAt     time.Time      `gorm:"index:idx_posts_author_created_id,priority:2"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	LikesCount    uint           `json:"likesCount" gorm:"not null;default:0"`
}

type Image struct {
	ID     uint   `gorm:"primaryKey"`
	PostID uint   `gorm:"not null;index"`
	URL    string `gorm:"not null"`
	Post   Post   `gorm:"constraint:OnDelete:CASCADE;"`
}

type AuthorView struct {
	ID        uint   `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

type PostView struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	ReplyToPostID *uint      `json:"replyToPostId,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	LikesCount    uint       `json:"likesCount"`
	Author        AuthorView `json:"author"`
}
