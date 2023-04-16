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

type HandleGetPlaylistRequest struct {
	ID         *int
	Identifier *string
}

func HandleGetPlaylist(w http.ResponseWriter, r *http.Request) {
	var req HandleGetPlaylistRequest

	// get the url vars
	params := mux.Vars(r)
	q := params["query"]

	// parse params
	if q != "" {
		intID, err := strconv.Atoi(q)
		if err != nil {
			// not an int, so it must be an identifier
			return
		}

		req.ID = &intID
	}

	if q != "" {
		req.Identifier = &q
	}

	// check if the user provided an id or an identifier
	if req.ID == nil && req.Identifier == nil {
		errors.ErrRequiredKey(w, "id or identifier")
		return
	}

	// execute the query
	playlist, err := query.GetPlaylist(db.DB, query.GetPlaylistRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
	})
	if err != nil {
		errors.SendError(w, "failed to execute query: "+err.Error(), http.StatusConflict)
		return
	}

	// send the response
	json.NewEncoder(w).Encode(responder.New(playlist))
}
