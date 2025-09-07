package handlers

import (
	"fmt"
	"net/http"
	"posts/internal/db/models"
	"posts/internal/dto"
	"posts/internal/feed"
	"posts/repo"
	"posts/utils"
	"strconv"
	"twitter/pkg/authjwt"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type PostsDeps struct {
	PostRepo  repo.Post
	Validator *validator.Validate
}

func GetPostById(deps *PostsDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		post, err := deps.PostRepo.GetPostById(req.Context(), uint(id))

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		utils.Res(w, http.StatusOK, post)
	}
}

func GetPosts(deps *PostsDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		cur, err := feed.DecodeCursor(r.URL.Query().Get("after"))
		if err != nil {
			http.Error(w, "bad cursor", http.StatusBadRequest)
			return
		}

		posts, next, err := deps.PostRepo.GetPosts(r.Context(), limit, cur)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		var nextStr *string
		if next != nil {
			s := feed.EncodeCursor(*next)
			nextStr = &s
		}

		utils.Res(w, http.StatusOK, map[string]any{
			"items": posts,
			"after": nextStr,
		})
	}
}

func CreatePost(deps *PostsDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var request dto.CreatePostRequest

		if err := utils.ParseJSON(w, req, &request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		err := deps.Validator.Struct(&request)

		if err != nil {
			http.Error(w, "invalid params", http.StatusBadRequest)
			return
		}

		id, ok := authjwt.UserIDFrom(req.Context())

		if !ok {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		post := &models.Post{
			Title:         request.Title,
			Description:   request.Description,
			AuthorID:      id,
			ReplyToPostID: request.ParentID,
		}

		err = deps.PostRepo.Create(req.Context(), post)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		message := fmt.Sprintf(`Успешно создан, id: %d`, post.ID)

		utils.Res(w, http.StatusCreated, utils.ResSuccess((message)))
	}
}

func UpdatePost(deps *PostsDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var request dto.UpdatePostRequest
		if err := utils.ParseAndValidateJSON(w, req, deps.Validator, &request); err != nil {
			http.Error(w, "invalid params", http.StatusBadRequest)
			return
		}

		userId, ok := authjwt.UserIDFrom(req.Context())

		if !ok {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		// TODO: не отправлять
		post, err := deps.PostRepo.GetPostById(req.Context(), request.Id)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if post == nil || userId != post.Author.ID {
			http.Error(w, "incorrect post id", http.StatusBadRequest)
			return
		}

		var updates = map[string]any{
			"title":       *request.Title,
			"description": *request.Description,
		}

		err = deps.PostRepo.Update(req.Context(), updates, post.ID)

		if err != nil {
			http.Error(w, "incorrect post id", http.StatusBadRequest)
			return
		}

		response := utils.ResSuccess(fmt.Sprintf(`Успешно отредактировано, ID: %d`, post.ID))

		utils.Res(w, http.StatusCreated, response)
	}
}

func DeletePost(deps *PostsDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var request dto.DeletePostRequest
		err := utils.ParseAndValidateJSON(w, req, deps.Validator, &request)

		if err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		userId, ok := authjwt.UserIDFrom(req.Context())

		if !ok {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		// TODO: не отправлять
		post, err := deps.PostRepo.GetPostById(req.Context(), request.Id)

		if err != nil || post.Author.ID != userId {
			http.Error(w, "invalid post id", http.StatusBadRequest)
		}

		err = deps.PostRepo.Delete(req.Context(), post.ID)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		response := utils.ResSuccess(fmt.Sprintf(`Успешно удалено, ID:%d`, post.ID))
		utils.Res(w, 200, response)
	}
}
