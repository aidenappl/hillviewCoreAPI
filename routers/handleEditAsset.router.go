package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
	"github.com/hillview.tv/coreAPI/structs"
)

type HandleEditAssetRequest struct {
	// Ident Fields from Query
	ID         *int    `json:"id"`
	Identifier *string `json:"identifier"`

	// Edit Fields
	Changes *assetChangeFields `json:"changes"`
}

type assetChangeFields struct {
	// Asset Fields
	Name        *string                `json:"name"`
	ImageURL    *string                `json:"image_url"`
	Identifier  *string                `json:"identifier"`
	Description *string                `json:"description"`
	Category    *int                   `json:"category"`
	Status      *int                   `json:"status"`
	Metadata    *structs.AssetMetadata `json:"metadata"`
}

func HandleEditAsset(w http.ResponseWriter, r *http.Request) {
	var req HandleEditAssetRequest
	// get mux variable
	q := mux.Vars(r)["query"]

	// check if query is valid
	if q != "" {
		intID, err := strconv.Atoi(q)
		if err != nil {
			// not an int, must be a string
			req.Identifier = &q
		}
		req.ID = &intID
	}

	// check that there is an identifier
	if req.ID == nil && req.Identifier == nil {
		responder.ErrRequiredKey(w, "id or identifier")
		return
	}

	// parse the body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responder.SendError(w, "failed to decode body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// check edit fields
	if req.Changes == nil {
		responder.ErrRequiredKey(w, "changes")
		return
	}

	if req.Changes.Name == nil && req.Changes.ImageURL == nil && req.Changes.Identifier == nil && req.Changes.Description == nil && req.Changes.Category == nil && req.Changes.Status == nil && req.Changes.Metadata == nil {
		responder.SendError(w, "no changes were provided", http.StatusBadRequest)
		return
	}

	// check if the asset exists
	asset, err := query.GetAsset(db.AssetDB, query.GetAssetRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
	})
	if err != nil {
		responder.SendError(w, "failed to get asset: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if asset == nil {
		responder.SendError(w, "asset not found", http.StatusNotFound)
		return
	}

	// update asset
	asset, err = query.UpdateAsset(db.AssetDB, query.UpdateAssetRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
		Changes: &query.UpdateAssetChanges{
			Name:        req.Changes.Name,
			ImageURL:    req.Changes.ImageURL,
			Identifier:  req.Changes.Identifier,
			Description: req.Changes.Description,
			Category:    req.Changes.Category,
			Status:      req.Changes.Status,
			Metadata:    req.Changes.Metadata,
		},
	})
	if err != nil {
		responder.SendError(w, "failed to update asset: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// return asset
	responder.New(w, asset)

}
