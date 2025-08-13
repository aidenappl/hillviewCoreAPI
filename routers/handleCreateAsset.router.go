package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"

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
		responder.SendError(w, "failed to decode body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// validate body
	if req.Name == nil || *req.Name == "" {
		responder.ParamError(w, "name")
		return
	}

	if req.Identifier == nil || *req.Identifier == "" {
		responder.ParamError(w, "identifier")
		return
	}

	if req.Category == nil {
		responder.ParamError(w, "category")
		return
	}

	if req.ImageURL == nil || *req.ImageURL == "" {
		responder.ParamError(w, "image_url")
		return
	}

	if req.Description == nil || *req.Description == "" {
		responder.ParamError(w, "description")
		return
	}

	// validate metadata
	if req.Metadata == nil {
		responder.ParamError(w, "metadata")
		return
	}

	if req.Metadata.Manufacturer == nil || *req.Metadata.Manufacturer == "" {
		responder.ParamError(w, "metadata.manufacturer")
		return
	}

	if req.Metadata.Model == nil || *req.Metadata.Model == "" {
		responder.ParamError(w, "metadata.model")
		return
	}

	if req.Metadata.SerialNumber == nil || *req.Metadata.SerialNumber == "" {
		responder.ParamError(w, "metadata.serial_number")
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
		responder.SendError(w, "failed to create asset: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responder.New(w, asset)

}
