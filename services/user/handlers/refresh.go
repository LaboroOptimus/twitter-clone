package handlers

import (
	"net/http"
	"strconv"
	"user/internal/dto"
	"user/utils"
)

func RefreshToken(deps LoginDeps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := utils.GetBearerToken(r.Header.Get("Authorization"))

		if err != nil {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := deps.JWT.ParseAndValidateRefresh(token)

		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		jti := claims.Jti
		userID64, _ := strconv.ParseUint(claims.Sub, 10, 64)
		userID := uint(userID64)

		u, err := deps.Users.GetByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		// получаем ua, remote addr
		ua, ip := r.UserAgent(), r.RemoteAddr

		// генерируем новый access
		accessToken, accessExp, err := deps.JWT.GenerateAccessToken(r.Context(), u.ID, u.Nickname)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// генерируем новый refresh
		newRefreshToken, newRefreshExp, newJTI, err := deps.JWT.GenerateRefreshToken(r.Context(), u.ID, u.Nickname)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		err = deps.Refresh.RevokeAndSave(r.Context(), jti, newJTI, userID, newRefreshExp, ua, ip)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		response := dto.LoginResponse{
			AccessToken: dto.TokenResponse{
				Token:     accessToken,
				ExpiresAt: accessExp,
			},
			RefreshToken: dto.TokenResponse{
				Token:     newRefreshToken,
				ExpiresAt: newRefreshExp,
			},
		}

		utils.Res(w, http.StatusOK, response)

	}
}
