package models

import "time"

type Follow struct {
	SubscriberID uint      `gorm:"not null;index:idx_sub_created_user,priority:1;uniqueIndex:ux_sub_user,priority:1"`
	UserID       uint      `gorm:"not null;index:idx_user_created_sub,priority:1;uniqueIndex:ux_sub_user,priority:2"`
	CreatedAt    time.Time `gorm:"not null;index:idx_sub_created_user,priority:2;index:idx_user_created_sub,priority:2"`

	Subscriber User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:SubscriberID;references:ID"`
	Followee   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
}
