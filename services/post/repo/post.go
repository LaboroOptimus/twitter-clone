package repo

import (
	"context"
	"posts/internal/db/models"

	"gorm.io/gorm"
)

type Post interface {
	GetPosts(ctx context.Context, page int) ([]*models.Post, error)
	GetPostById(ctx context.Context, id uint) (*models.Post, error)
	Create(ctx context.Context, p *models.Post) error
}

type post struct{ db *gorm.DB }

func NewPost(db *gorm.DB) *post {
	return &post{db: db}
}

func (p *post) Create(ctx context.Context, newPost *models.Post) error {
	sql := `
	  INSERT INTO posts (title, description, reply_to_post_id, author_id, created_at)
	  VALUES (?, ?, ?, ?, ?)
	  RETURNING id, created_at`

	row := p.db.WithContext(ctx).Raw(sql,
		newPost.Title,
		newPost.Description,
		newPost.ReplyToPostID,
		newPost.AuthorID,
		newPost.CreatedAt,
	).Row()

	if err := row.Scan(&newPost.ID, &newPost.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (p *post) GetPosts(ctx context.Context, page int) ([]*models.Post, error) {
	const limit = 20
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var posts []*models.Post

	rows := p.db.WithContext(ctx).
		Raw(`SELECT id, title, description, reply_to_post_id, author_id, created_at, updated_at
         FROM posts
         ORDER BY created_at DESC
         LIMIT ? OFFSET ?`, limit, offset)

	if err := rows.Scan(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *post) GetPostById(ctx context.Context, id uint) (*models.Post, error) {
	sql := `SELECT id, title, description, reply_to_post_id, author_id, created_at, updated_at
		FROM posts
		WHERE id = ?
	`

	row := p.db.WithContext(ctx).Raw(sql, id)

	var Post *models.Post

	if err := row.Scan(&Post).Error; err != nil {
		return nil, err
	}

	return Post, nil
}
