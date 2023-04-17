package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/errors"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleCreateAssetRequest struct {
	query.CreateAssetRequest
}

func HandleCreateAsset(w http.ResponseWriter, r *http.Request) {
	var req HandleCreateAssetRequest

	// decode body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.SendError(w, "failed to decode body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// validate body
	if req.Name == nil || *req.Name == "" {
		errors.ParamError(w, "name")
		return
	}

	if req.Identifier == nil || *req.Identifier == "" {
		errors.ParamError(w, "identifier")
		return
	}

	if req.Category == nil {
		errors.ParamError(w, "category")
		return
	}

	if req.ImageURL == nil || *req.ImageURL == "" {
		errors.ParamError(w, "image_url")
		return
	}

	if req.Description == nil || *req.Description == "" {
		errors.ParamError(w, "description")
		return
	}

	// validate metadata
	if req.Metadata == nil {
		errors.ParamError(w, "metadata")
		return
	}

	if req.Metadata.Manufacturer == nil || *req.Metadata.Manufacturer == "" {
		errors.ParamError(w, "metadata.manufacturer")
		return
	}

	if req.Metadata.Model == nil || *req.Metadata.Model == "" {
		errors.ParamError(w, "metadata.model")
		return
	}

	if req.Metadata.SerialNumber == nil || *req.Metadata.SerialNumber == "" {
		errors.ParamError(w, "metadata.serial_number")
		return
	}

	// create asset
	asset, err := query.CreateAsset(db.AssetDB, query.CreateAssetRequest{
		Name:        req.Name,
		ImageURL:    req.ImageURL,
		Identifier:  req.Identifier,
		Description: req.Description,
		Category:    req.Category,
		Metadata:    req.Metadata,
	})
	if err != nil {
		errors.SendError(w, "failed to create asset: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(responder.New(asset))

}
