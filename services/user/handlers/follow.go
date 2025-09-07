package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"twitter/pkg/authjwt"
	"twitter/pkg/redisx"
	"user/repo"
	"user/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type FollowDeps struct {
	Follow      repo.Follow
	Producer    *redisx.Producer
	NotifStream string
}

func FollowUser(deps *FollowDeps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)

		if err != nil || id < 0 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		userID, ok := authjwt.UserIDFrom(req.Context())

		if !ok {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if userID == uint(id) {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		err = deps.Follow.Follow(req.Context(), userID, uint(id))
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		pubCtx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		defer cancel()

		_, _ = deps.Producer.Publish(pubCtx, deps.NotifStream, map[string]interface{}{
			"type":       "follow",
			"user_id":    strconv.FormatUint(uint64(id), 10), 
			"actor_id":   strconv.FormatUint(uint64(userID), 10),
			"event_uuid": uuid.NewString(),
		})

		utils.Res(w, http.StatusOK, nil)
	}
}

func UnfollowUser(deps *FollowDeps) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)

		if err != nil || id < 0 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		userID, ok := authjwt.UserIDFrom(req.Context())

		if !ok {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		if userID == uint(id) {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		err = deps.Follow.Unfollow(req.Context(), userID, uint(id))

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		utils.Res(w, http.StatusOK, nil)

	}
}
