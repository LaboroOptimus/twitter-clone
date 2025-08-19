package handlers

import (
	"errors"
	"net/http"
	"user/internal/db/models"
	"user/internal/dto"
	"user/internal/validation"
	"user/repo"
	"user/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

func Register(userRepo repo.User) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request dto.RegisterRequest

		if err := utils.ParseJSON(w, r, &request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if err := validation.Validate.Struct(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := models.NewUser(request.Nickname, request.Email, request.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := userRepo.Create(r.Context(), user); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				http.Error(w, "email or nickname already exists", http.StatusConflict)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		response := utils.GenerateResponse("User created successfully", utils.Success)

		utils.Res(w, http.StatusCreated, response)
	}
}
