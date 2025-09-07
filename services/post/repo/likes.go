package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Likes interface {
	Like(ctx context.Context, userID uint, postID uint) error
	Unlike(ctx context.Context, userID uint, postID uint) error
	LikesAmount(ctx context.Context, postID uint) (*int64, error)
}

type like struct{ db *gorm.DB }

func NewLike(db *gorm.DB) *like {
	return &like{db: db}
}

func (l *like) Like(ctx context.Context, userID, postID uint) error {
	return l.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(`
      INSERT INTO likes (user_id, post_id, created_at)
      VALUES (?, ?, ?)
      ON CONFLICT (post_id, user_id) DO NOTHING
    `, userID, postID, time.Now())
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 1 {
			upd := tx.Exec(`
        UPDATE posts SET likes_count = likes_count + 1
        WHERE id = ?
      `, postID)
			if upd.Error != nil {
				return upd.Error
			}
		}
		return nil
	})
}
func (l *like) Unlike(ctx context.Context, userID uint, postID uint) error {
	return l.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Exec(`
		DELETE FROM likes 
		WHERE user_id = ? AND post_id = ?
		`, userID, postID)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 1 {
			upd := tx.Exec(`
        UPDATE posts SET likes_count = likes_count - 1
        WHERE id = ?
      `, postID)
			if upd.Error != nil {
				return upd.Error
			}
		}
		return nil
	})
}

func (l *like) LikesAmount(ctx context.Context, postID uint) (*int64, error) {
	var likesCount int64
	tx := l.db.WithContext(ctx).
		Raw(`SELECT likes_count FROM posts WHERE id = ?`, postID).
		Scan(&likesCount)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &likesCount, nil
}
