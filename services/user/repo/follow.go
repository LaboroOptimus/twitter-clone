package repo

import (
	"context"
	"errors"
	"user/internal/db/models"

	"gorm.io/gorm"
)

type Follow interface {
	Follow(ctx context.Context, userID uint, followID uint) error
	Unfollow(ctx context.Context, userID uint, followID uint) error
}

type follow struct {
	db *gorm.DB
}

func NewFollow(db *gorm.DB) *follow {
	return &follow{
		db: db,
	}
}

func (f *follow) Follow(ctx context.Context, userID uint, followID uint) error {
	err := f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(`
            INSERT INTO follows (user_id, subscriber_id, created_at)
            VALUES (?, ?, NOW())
            ON CONFLICT (user_id, subscriber_id) DO NOTHING
        `, userID, followID)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return errors.New("already subscribed")
		}

		if res.RowsAffected == 1 {
			if err := tx.Model(&models.User{}).
				Where("id = ?", userID).
				UpdateColumn("subscribers_amount", gorm.Expr("subscribers_amount + 1")).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (f *follow) Unfollow(ctx context.Context, userID uint, followID uint) error {

	err := f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(`
			DELETE FROM follows
			WHERE user_id = ? AND subscriber_id = ?
		`, userID, followID)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return errors.New("already unsubscribed")
		}


		if res.RowsAffected == 1 {
			if err := tx.Model(&models.User{}).
				Where("id = ? AND subscribers_amount > 0", userID).
				UpdateColumn("subscribers_amount", gorm.Expr("subscribers_amount - 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
