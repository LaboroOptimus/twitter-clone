package repo

import (
	"context"
	"errors"
	"time"
	"user/internal/db/models"

	"gorm.io/gorm"
)

type Refresh interface {
	Save(ctx context.Context, userID uint, jti string, exp time.Time, ua, ip string) error
	Revoke(ctx context.Context, jti string) (bool, error)
	IsRevokedOrExpired(ctx context.Context, jti string) (bool, error)
	RevokeAndSave(ctx context.Context, oldJTI string, newJTI string, userID uint, exp time.Time, ua, ip string) error
}

type refresh struct{ db *gorm.DB }

func NewRefresh(db *gorm.DB) Refresh { return &refresh{db: db} }

func (r *refresh) Save(ctx context.Context, userID uint, jti string, exp time.Time, ua, ip string) error {
	refresh := models.Refresh{
		UserID:    userID,
		JTI:       jti,
		ExpiresAt: exp,
		UserAgent: ua,
		IP:        ip,
	}

	return r.db.WithContext(ctx).Create(&refresh).Error
}
func (r *refresh) Revoke(ctx context.Context, jti string) (bool, error) {
	now := time.Now()
	tx := r.db.WithContext(ctx).Model(&models.Refresh{}).
		Where("jti = ? AND revoked_at IS NULL", jti).
		Update("revoked_at", now)
	return tx.RowsAffected > 0, tx.Error
}

func (r *refresh) IsRevokedOrExpired(ctx context.Context, jti string) (bool, error) {
	var rec models.Refresh
	err := r.db.WithContext(ctx).
		Select("expires_at", "revoked_at").
		Where("jti = ?", jti).
		Take(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}
	if rec.RevokedAt != nil {
		return true, nil
	}
	if time.Now().After(rec.ExpiresAt) {
		return true, nil
	}
	return false, nil
}

func (r *refresh) RevokeAndSave(ctx context.Context, oldJTI string, newJTI string, userID uint, exp time.Time, ua, ip string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		t := tx.Model(&models.Refresh{}).Where("jti = ? AND revoked_at IS NULL", oldJTI).Update("revoked_at", now)

		if t.Error != nil || t.RowsAffected == 0 {
			return t.Error
		}

		refresh := models.Refresh{
			UserID:    userID,
			JTI:       newJTI,
			ExpiresAt: exp,
			UserAgent: ua,
			IP:        ip,
		}

		err := tx.Create(&refresh).Error

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
