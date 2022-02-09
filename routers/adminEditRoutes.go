package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
)

type HandleEditAssetRequest struct {
	ID          *int    `json:"id"`
	Name        *string `json:"name"`
	ImageURL    *string `json:"image_url"`
	Identifier  *string `json:"identifier"`
	Description *string `json:"description"`
	Category    *int    `json:"category"`
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
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(asset)
}
