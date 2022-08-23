package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
)

type HandleEditMobileAccountRequest struct {
	ID     *int                              `json:"id"`
	Fields *query.MobileAccountModifications `json:"fields"`
}

func HandleEditMobileAccount(w http.ResponseWriter, r *http.Request) {
	body := HandleEditMobileAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.ID == nil || body.Fields == nil {
		http.Error(w, "missing required keys", http.StatusBadRequest)
		return
	}

	user, err := query.EditMobileAccount(db.AssetDB, query.EditMobileAccountRequest{
		ID:            body.ID,
		Modifications: body.Fields,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

type HandleEditAssetRequest struct {
	ID          *int    `json:"id"`
	Name        *string `json:"name"`
	ImageURL    *string `json:"image_url"`
	Identifier  *string `json:"identifier"`
	Description *string `json:"description"`
	Category    *int    `json:"category"`
	Notes       *string `json:"notes"`
	Status      *int    `json:"status"`
}

func HandleEditAsset(w http.ResponseWriter, r *http.Request) {
	body := HandleEditAssetRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if body.ID == nil && !(body.Name == nil && body.ImageURL == nil && body.Identifier == nil && body.Description == nil && body.Category == nil && body.Status == nil) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	asset, err := query.EditAsset(db.AssetDB, query.EditAssetRequest{
		ID: body.ID,
		Modifications: &query.AssetModifications{
			Name:        body.Name,
			ImageURL:    body.ImageURL,
			Identifier:  body.Identifier,
			Description: body.Description,
			Category:    body.Category,
			Status:      body.Status,
			Notes:       body.Notes,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(asset)
}

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
