package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/errors"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleCreatePlaylistRequest struct {
	query.CreatePlaylistRequest
}

func HandleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req HandleCreatePlaylistRequest

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

	if req.Description == nil || *req.Description == "" {
		errors.ParamError(w, "description")
		return
	}

	if req.BannerImage == nil || *req.BannerImage == "" {
		errors.ParamError(w, "banner_image")
		return
	}

	if req.Route == nil || *req.Route == "" {
		errors.ParamError(w, "route")
		return
	}

	if req.Videos == nil || len(*req.Videos) == 0 {
		errors.ParamError(w, "videos")
		return
	}

	// check videos
	for _, videoID := range *req.Videos {
		// get video
		video, err := query.GetVideo(db.DB, query.GetVideoRequest{
			ID: &videoID,
		})
		if err != nil {
			errors.SendError(w, "failed to get video: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// check if video exists
		if video == nil {
			errors.SendError(w, fmt.Sprintf("video: %d, does not exist", videoID), http.StatusBadRequest)
			return
		}
	}

	// create playlist
	playlist, err := query.CreatePlaylist(db.DB, req.CreatePlaylistRequest)
	if err != nil {
		errors.SendError(w, "failed to create playlist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send response
	json.NewEncoder(w).Encode(responder.New(playlist))

}
