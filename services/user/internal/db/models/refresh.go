package models

import "time"

type Refresh struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID"`
	JTI       string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	UserAgent string
	IP        string
	CreatedAt time.Time
}

func NewRefresh(userID uint, jti string, expiresAt time.Time, ua, ip string) (*Refresh, error) {
	return &Refresh{
		UserID:    userID,
		JTI:       jti,
		ExpiresAt: expiresAt,
		RevokedAt: nil,
		UserAgent: ua,
		IP:        ip,
	}, nil
}
