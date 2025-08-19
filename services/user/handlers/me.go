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

func GetProfile(userRepo repo.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := r.FormValue("id")

		if idParam == "" {
			http.Error(w, "invalid request", http.StatusBadRequest)
		}

		id, err := strconv.Atoi(idParam)

		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
		}

		user, err := userRepo.GetByID(r.Context(), uint(id))

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		result := dto.ProfileResponse{
			ID:        user.ID,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Bio:       user.Bio,
			AvatarURL: user.AvatarURL,
		}

		utils.Res(w, http.StatusOK, result)

	}
}

func UpdateProfile(userRepo repo.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := authjwt.UserIDFrom(r.Context())

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

		if err := userRepo.UpdateProfile(r.Context(), userID, updates); err != nil {
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
