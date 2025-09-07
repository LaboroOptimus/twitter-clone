package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"twitter/pkg/authjwt"
	"user/internal/dto"
	"user/internal/validation"
	"user/repo"
	"user/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserDeps struct {
	Users repo.User
}

func GetProfile(deps UserDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := authjwt.UserIDFrom(r.Context())

		if !ok {
			http.Error(w, "invalid request", http.StatusBadRequest)
		}

		user, err := deps.Users.GetByID(r.Context(), id)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		result := dto.ProfileResponse{
			ID:                user.ID,
			Email:             user.Email,
			Nickname:          user.Nickname,
			Bio:               user.Bio,
			AvatarURL:         user.AvatarURL,
			SubscribersAmount: user.SubscribersAmount,
		}

		utils.Res(w, http.StatusOK, result)

	}
}

func UpdateProfile(deps UserDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := authjwt.UserIDFrom(r.Context())

		if !ok {
			http.Error(w, "invalid request", http.StatusBadRequest)
		}

		var request dto.UpdateRequest
		if err := utils.ParseJSON(w, r, &request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if err := validation.Validate.Struct(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// собираем только пришедшие поля
		updates := map[string]interface{}{}
		if request.Nickname != "" {
			updates["nickname"] = request.Nickname
		}
		if request.Bio != "" {
			updates["bio"] = request.Bio
		}
		if request.AvatarURL != "" {
			updates["avatar_url"] = request.AvatarURL
		}

		if err := deps.Users.UpdateProfile(r.Context(), userID, updates); err != nil {
			// ловим конфликт уникальности никнейма (23505)
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				http.Error(w, "nickname already taken", http.StatusConflict)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		utils.Res(w, http.StatusOK, map[string]string{"message": "profile updated"})
	}
}

func GetUser(deps UserDeps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		idParam := req.URL.Query().Get("id")

		id, err := strconv.Atoi(idParam)
		if err != nil || id < 0 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		user, err := deps.Users.GetByID(req.Context(), uint(id))

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		response := dto.ProfileResponse{
			ID:                user.ID,
			Email:             user.Email,
			Nickname:          user.Nickname,
			Bio:               user.Bio,
			AvatarURL:         user.AvatarURL,
			SubscribersAmount: user.SubscribersAmount,
		}

		utils.Res(w, http.StatusOK, response)
	}
}
