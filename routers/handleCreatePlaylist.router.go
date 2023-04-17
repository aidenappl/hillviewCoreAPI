package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/errors"
	"github.com/hillview.tv/coreAPI/query"
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
	
}
