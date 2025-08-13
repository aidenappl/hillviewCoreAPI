package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleUpdatePlaylistRequest struct {
	ID      *int                         `json:"playlist_id"`
	Changes *query.UpdatePlaylistChanges `json:"changes"`
}

func HandleUpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	var req HandleUpdatePlaylistRequest

	// get from query params
	q := mux.Vars(r)["query"]

	// parse query params
	if q != "" {
		intID, err := strconv.Atoi(q)
		if err != nil {
			responder.SendError(w, "failed to parse query param", http.StatusBadRequest)
			return
		}
		req.ID = &intID
	} else {
		responder.ParamError(w, "query")
		return
	}

	// get from body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responder.SendError(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	// validate body
	if req.Changes == nil {
		responder.ParamError(w, "changes")
		return
	}

	// update playlist
	playlist, err := query.UpdatePlaylist(db.DB, query.UpdatePlaylistRequest{
		ID:      req.ID,
		Changes: req.Changes,
	})
	if err != nil {
		responder.SendError(w, "failed to update playlist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send response
	responder.New(w, playlist)
}
