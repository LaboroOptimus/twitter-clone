package repo

import (
	"context"
	"posts/internal/db/models"
	"posts/internal/feed"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Post interface {
	GetPosts(ctx context.Context, limit int, after *feed.Cursor) ([]models.PostView, *feed.Cursor, error)
	GetPostById(ctx context.Context, id uint) (*models.PostView, error)
	Create(ctx context.Context, p *models.Post) error
	Update(ctx context.Context, updates map[string]any, postID uint) error
	Delete(ctx context.Context, postID uint) error
}

type post struct{ db *gorm.DB }

func NewPost(db *gorm.DB) *post {
	return &post{db: db}
}

type postRow struct {
	ID              uint
	Title           string
	Description     string
	ReplyToPostID   *uint
	CreatedAt       time.Time
	LikesCount      uint
	AuthorID        uint   `gorm:"column:author_id"`
	AuthorNickname  string `gorm:"column:author_nickname"`
	AuthorAvatarURL string `gorm:"column:author_avatar_url"`
}

func (p *post) Create(ctx context.Context, newPost *models.Post) error {
	sql := `
	  INSERT INTO posts (title, description, reply_to_post_id, author_id, created_at, updated_at)
	  VALUES (?, ?, ?, ?, ?, ?)
	  RETURNING id, created_at`

	if err := p.db.WithContext(ctx).Raw(sql,
		newPost.Title,
		newPost.Description,
		newPost.ReplyToPostID,
		newPost.AuthorID,
		time.Now(),
		time.Now(),
	).Scan(newPost).Error; err != nil {
		return err
	}
	return nil
}

func (p *post) GetPosts(ctx context.Context, limit int, after *feed.Cursor) ([]models.PostView, *feed.Cursor, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	q := p.db.WithContext(ctx).
		Table("posts p").
		Select("p.id, p.title, p.description, p.reply_to_post_id, p.author_id, p.created_at, p.updated_at, p.likes_count").
		Joins("JOIN users u ON u.id = p.author_id").
		Where("p.deleted_at IS NULL")

	if after != nil {
		q = q.Where("(created_at < ?) OR (created_at = ? AND id < ?)", after.CreatedAt, after.CreatedAt, after.ID)
	}

	var rows []postRow

	if err := q.Order("p.created_at DESC, p.id DESC").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, nil, err
	}

	out := make([]models.PostView, 0, len(rows))

	for _, r := range rows {
		out = append(out, models.PostView{
			ID: r.ID, Title: r.Title, Description: r.Description,
			ReplyToPostID: r.ReplyToPostID, CreatedAt: r.CreatedAt, LikesCount: r.LikesCount,
			Author: models.AuthorView{ID: r.AuthorID, Nickname: r.AuthorNickname, AvatarURL: r.AuthorAvatarURL},
		})
	}

	var next *feed.Cursor
	if len(rows) == limit {
		last := rows[len(rows)-1]
		next = &feed.Cursor{CreatedAt: last.CreatedAt, ID: last.ID}
	}
	return out, next, nil
}

func (p *post) GetPostById(ctx context.Context, id uint) (*models.PostView, error) {

	var r postRow
	tx := p.db.WithContext(ctx).
		Table("posts AS p").
		Select(`
			p.id, p.title, p.description, p.reply_to_post_id,
			p.author_id, p.created_at, p.updated_at, p.likes_count,
			u.nickname AS author_nickname, u.avatar_url AS author_avatar_url
		`).
		Joins("INNER JOIN users AS u ON u.id = p.author_id").
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		Limit(1).
		Scan(&r)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	out := &models.PostView{
		ID:            r.ID,
		Title:         r.Title,
		Description:   r.Description,
		ReplyToPostID: r.ReplyToPostID,
		CreatedAt:     r.CreatedAt,
		LikesCount:    r.LikesCount,
		Author: models.AuthorView{
			ID:        r.AuthorID,
			Nickname:  r.AuthorNickname,
			AvatarURL: r.AuthorAvatarURL,
		},
	}
	return out, nil
}

func (p *post) Update(ctx context.Context, updates map[string]any, postID uint) error {
	allowed := map[string]bool{"title": true, "description": true}
	parts := make([]string, 0, len(updates)+1)
	args := make([]any, 0, len(updates)+2)

	for k, v := range updates {
		if !allowed[k] {
			continue
		}
		parts = append(parts, k+" = ?")
		args = append(args, v)
	}
	parts = append(parts, "updated_at = now()")

	query := "UPDATE posts SET " + strings.Join(parts, ", ") + " WHERE id = ?"
	args = append(args, postID)

	tx := p.db.WithContext(ctx).Exec(query, args...)

	if tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (p *post) Delete(ctx context.Context, id uint) error {
	tx := p.db.WithContext(ctx).Exec(`UPDATE posts SET deleted_at = now() WHERE id = ?`, id)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
