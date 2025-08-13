package routers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleGetLinkRequest struct {
	ID         *int
	Identifier *string
}

func HandleGetLink(w http.ResponseWriter, r *http.Request) {

	var req HandleGetLinkRequest

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
		responder.ErrRequiredKey(w, "id or identifier")
		return
	}

	// execute the query
	link, err := query.GetLink(db.DB, query.GetLinkRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
	})
	if err != nil {
		responder.SendError(w, "failed to get videos: "+err.Error(), http.StatusConflict)
		return
	}

	// send the response
	responder.New(w, link)

}
