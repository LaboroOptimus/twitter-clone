package dto

type ProfileResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatarURL"`
	SubscribersAmount uint `json:"subscribersAmount"`
}

type UserProfileRequest struct {
	ID uint `json:"userId"`
}
