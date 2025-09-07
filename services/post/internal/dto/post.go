package dto

type CreatePostRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=1,max=300"`
	Image       string `json:"image" validate:"omitempty,url"`
	ParentID    *uint  `json:"parentId,omitempty" validate:"omitempty,gt=0"`
}

type UpdatePostRequest struct {
	Id          uint    `json:"id" validate:"required"`
	Title       *string `json:"title" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description" validate:"omitempty,min=1,max=300"`
}

type DeletePostRequest struct {
	Id uint `json:"id" validate:"required"`
}
