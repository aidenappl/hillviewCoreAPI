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

type HandleGetVideoRequest struct {
	ID         *int
	Identifier *string
}

func HandleGetVideo(w http.ResponseWriter, r *http.Request) {

	var req HandleGetVideoRequest

	// get the url vars
	params := mux.Vars(r)
	q := params["query"]

	// parse params
	if q != "" {
		intID, err := strconv.Atoi(q)
		if err != nil {
			req.Identifier = &q
		} else {
			req.ID = &intID
		}
	}

	// check if the user provided an id or an identifier
	if req.ID == nil && req.Identifier == nil {
		errors.ErrRequiredKey(w, "id or identifier")
		return
	}

	// execute the query
	videos, err := query.GetVideo(db.DB, query.GetVideoRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
	})
	if err != nil {
		errors.SendError(w, "failed to get videos: "+err.Error(), http.StatusConflict)
		return
	}

	// send the response
	json.NewEncoder(w).Encode(responder.New(videos))

}
