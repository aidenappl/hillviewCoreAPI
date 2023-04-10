package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
)

type HandleEditVideoRequest struct {
	ID          *int    `json:"id"`
	Title       *string `json:"title"`
	Thumbnail   *string `json:"thumbnail"`
	Description *string `json:"description"`
	URL         *string `json:"url"`
	Status      *int    `json:"status"`
}

func HandleEditVideo(w http.ResponseWriter, r *http.Request) {
	body := HandleEditVideoRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.ID == nil && !(body.Title == nil && body.Thumbnail == nil && body.Description == nil && body.URL == nil) {
		http.Error(w, "missing required keys", http.StatusBadRequest)
		return
	}

	asset, err := query.EditVideo(db.DB, query.EditVideoRequest{
		ID: body.ID,
		Modifications: &query.VideoModifications{
			Title:       body.Title,
			Thumbnail:   body.Thumbnail,
			Description: body.Description,
			URL:         body.URL,
			Status:      body.Status,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(asset)
}

type HandleEditAdminRequest struct {
	ID             *int    `json:"id"`
	Name           *string `json:"name"`
	Email          *string `json:"email"`
	Username       *string `json:"username"`
	Authentication *int    `json:"authentication"`
}

func HandleEditAdminAccount(w http.ResponseWriter, r *http.Request) {
	body := HandleEditAdminRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.ID == nil && !(body.Name == nil && body.Email == nil && body.Authentication == nil && body.Username == nil) {
		http.Error(w, "missing required keys", http.StatusBadRequest)
		return
	}

	user, err := query.EditAdminAccount(db.DB, query.EditAdminAccountRequest{
		ID: *body.ID,
		Modifications: query.AdminAccountModifications{
			Name:           body.Name,
			Email:          body.Email,
			Authentication: body.Authentication,
			Username:       body.Username,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
