package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/errors"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleCreateVideoRequest struct {
	query.CreateVideoRequest
}

func HandleCreateVideo(w http.ResponseWriter, r *http.Request) {
	var req HandleCreateVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.SendError(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// validate the request
	if req.Title == nil || *req.Title == "" {
		errors.ParamError(w, "title")
		return
	}

	if req.Description == nil || *req.Description == "" {
		errors.ParamError(w, "description")
		return
	}

	if req.Thumbnail == nil || *req.Thumbnail == "" {
		errors.ParamError(w, "thumbnail")
		return
	}

	if req.URL == nil || *req.URL == "" {
		errors.ParamError(w, "url")
		return
	}

	if req.Status == nil {
		errors.ParamError(w, "status")
		return
	}

	// create the video
	video, err := query.CreateVideo(db.DB, req.CreateVideoRequest)
	if err != nil {
		errors.SendError(w, "failed to create video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send the response
	json.NewEncoder(w).Encode(responder.New(video))

}
