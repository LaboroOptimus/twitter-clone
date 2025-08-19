package authjwt

import (
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware(ts TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ah := r.Header.Get("Authorization")
			if !strings.HasPrefix(ah, "Bearer ") {
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(ah, "Bearer ")

			claims, err := ts.ParseAndValidateAccess(tokenStr)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			uid64, _ := strconv.ParseUint(claims.Sub, 10, 64)
			ctx := WithUserID(r.Context(), uint(uid64))
			ctx = WithNickname(ctx, claims.Subject)
			ctx = WithJTI(ctx, claims.Jti)

			if uid64 == 0 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
