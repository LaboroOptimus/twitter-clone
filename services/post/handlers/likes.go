package handlers

import (
	"context"
	"net/http"
	"posts/internal/dto"
	"posts/repo"
	"posts/utils"
	"time"
	"twitter/pkg/authjwt"
	"twitter/pkg/redisx"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type LikesDeps struct {
	LikeRepo    repo.Likes
	PostRepo    repo.Post
	Validator   *validator.Validate
	Producer    *redisx.Producer
	NotifStream string
}

func LikePost(deps *LikesDeps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var request dto.LikeUnlikePostRequest

		err := utils.ParseAndValidateJSON(w, req, deps.Validator, &request)

		if err != nil {
			http.Error(w, "invalid param postId", http.StatusBadRequest)
			return
		}

		// TODO: - надо ли вообще это делать? проще разрулить на уровне бд
		post, err := deps.PostRepo.GetPostById(req.Context(), request.PostID)
		if err != nil {
			http.Error(w, "invalid param postId", http.StatusBadRequest)
			return
		}

		userId, ok := authjwt.UserIDFrom(req.Context())

		if !ok {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err = deps.LikeRepo.Like(req.Context(), userId, request.PostID)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		pubCtx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		defer cancel()

		_, _ = deps.Producer.Publish(pubCtx, deps.NotifStream, map[string]interface{}{
			"type":       "like",
			"user_id":    post.Author.ID,
			"actor_id":   userId,
			"post_id":    request.PostID,
			"event_uuid": uuid.NewString(),
		})

		response := utils.ResSuccess("")
		utils.Res(w, 200, response)
	}
}

func UnlikePost(deps *LikesDeps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var request dto.LikeUnlikePostRequest

		if err := utils.ParseAndValidateJSON(w, req, deps.Validator, &request); err != nil {
			http.Error(w, "invalid params", http.StatusBadRequest)
			return
		}

		// TODO: - надо ли вообще это делать? проще разрулить на уровне бд
		_, err := deps.PostRepo.GetPostById(req.Context(), request.PostID)
		if err != nil {
			http.Error(w, "invalid param postId", http.StatusBadRequest)
			return
		}

		userId, ok := authjwt.UserIDFrom(req.Context())

		if !ok {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err = deps.LikeRepo.Unlike(req.Context(), userId, request.PostID)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		response := utils.ResSuccess("")
		utils.Res(w, 200, response)

	}
}
