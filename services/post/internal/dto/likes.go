package dto

type LikeUnlikePostRequest struct {
	PostID uint `json:"postId" validate:"required"`
}
