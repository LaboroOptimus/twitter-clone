package handlers

import (
	"net/http"
	"twitter/pkg/authjwt"
	"user/internal/auth"
	"user/internal/dto"
	"user/internal/validation"
	"user/repo"
	"user/utils"
)

type LoginDeps struct {
	Users   repo.User
	Refresh repo.Refresh
	JWT     authjwt.TokenService
}

func Login(deps LoginDeps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request dto.LoginRequest

		if err := utils.ParseJSON(w, r, &request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if err := validation.Validate.Struct(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := deps.Users.GetByEmail(r.Context(), request.Email)

		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		checked := auth.Compare(user.Password, request.Password)

		if !checked {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		accessToken, expAccess, err := deps.JWT.GenerateAccessToken(r.Context(), user.ID, user.Nickname)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		refreshToken, expRefresh, jti, err := deps.JWT.GenerateRefreshToken(r.Context(), user.ID, user.Nickname)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		ua := r.UserAgent()
		ip := r.RemoteAddr

		err = deps.Refresh.Save(r.Context(), user.ID, jti, expRefresh, ua, ip)

		if err != nil {
			http.Error(w, "failed to save refresh token", http.StatusInternalServerError)
			return
		}

		response := dto.LoginResponse{
			AccessToken: dto.TokenResponse{
				Token:     accessToken,
				ExpiresAt: expAccess,
			},
			RefreshToken: dto.TokenResponse{
				Token:     refreshToken,
				ExpiresAt: expRefresh,
			},
		}

		utils.Res(w, http.StatusOK, response)
	}
}
