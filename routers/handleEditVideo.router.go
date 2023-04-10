package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/errors"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleEditVideoRequest struct {
	// Ident Fields from Query
	ID         *int    `json:"id"`
	Identifier *string `json:"identifier"`

	// Edit Fields
	Changes *videoChangeFields `json:"changes"`
}

type videoChangeFields struct {
	// Video Fields
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	Thumbnail      *string `json:"thumbnail"`
	AllowDownloads *bool   `json:"allow_downloads"`
	DownloadURL    *string `json:"download_url"`
	URL            *string `json:"url"`
	Status         *int    `json:"status"`
}

func HandleEditVideo(w http.ResponseWriter, r *http.Request) {
	var req HandleEditVideoRequest
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
		errors.ErrRequiredKey(w, "id or identifier")
		return
	}

	// parse the body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.SendError(w, "failed to decode body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// check edit fields
	if req.Changes == nil {
		errors.ErrRequiredKey(w, "changes")
		return
	}

	if req.Changes.Title == nil && req.Changes.Description == nil && req.Changes.Thumbnail == nil && req.Changes.AllowDownloads == nil && req.Changes.DownloadURL == nil && req.Changes.URL == nil && req.Changes.Status == nil {
		errors.SendError(w, "no changes to make", http.StatusBadRequest)
		return
	}

	// check if the asset exists
	video, err := query.GetVideo(db.DB, query.GetVideoRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
	})
	if err != nil {
		errors.SendError(w, "failed to get video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if video == nil {
		errors.SendError(w, "video not found", http.StatusNotFound)
		return
	}

	// update asset
	video, err = query.UpdateVideo(db.DB, query.UpdateVideoRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
		Changes: &query.UpdateVideoChanges{
			Title:          req.Changes.Title,
			Description:    req.Changes.Description,
			Thumbnail:      req.Changes.Thumbnail,
			AllowDownloads: req.Changes.AllowDownloads,
			DownloadURL:    req.Changes.DownloadURL,
			URL:            req.Changes.URL,
			Status:         req.Changes.Status,
		},
	})
	if err != nil {
		errors.SendError(w, "failed to update video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// return asset
	json.NewEncoder(w).Encode(responder.New(video))
}
