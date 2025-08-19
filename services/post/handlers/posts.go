package handlers

import (
	"fmt"
	"net/http"
	"posts/repo"
	"strconv"
)

func GetPosts(postRepo repo.Post) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		page := req.FormValue("page")
		pageParam := 0

		if page == "" {
			pageParam = -1
		} else {
			val, err := strconv.Atoi(page)

			if err != nil {
				http.Error(w, "invalid page param", http.StatusBadRequest)
			}
			pageParam = val
		}

		// var posts []models.Post;
		posts, err := postRepo.GetPosts(req.Context(), pageParam)

		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		fmt.Println(posts)

	}
}
