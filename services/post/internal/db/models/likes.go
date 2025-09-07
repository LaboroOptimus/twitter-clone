package models

import "time"

type Likes struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	PostID    uint `gorm:"not null;uniqueIndex:idx_post_user,priority:1"`
	UserID    uint `gorm:"not null;uniqueIndex:idx_post_user,priority:2"`
	CreatedAt time.Time
}
