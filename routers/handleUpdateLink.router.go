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

type HandleUpdateLinkRequest struct {
	ID      *int                     `json:"link_id"`
	Changes *query.UpdateLinkChanges `json:"changes"`
}

func HandleUpdateLink(w http.ResponseWriter, r *http.Request) {
	var req HandleUpdateLinkRequest

	// get from query params
	q := mux.Vars(r)["query"]

	// parse query params
	if q != "" {
		intID, err := strconv.Atoi(q)
		if err != nil {
			errors.SendError(w, "failed to parse query param", http.StatusBadRequest)
			return
		}
		req.ID = &intID
	} else {
		errors.ParamError(w, "query")
		return
	}

	// get from body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.SendError(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	// validate body
	if req.Changes == nil {
		errors.ParamError(w, "changes")
		return
	}

	// update link
	link, err := query.UpdateLink(db.DB, query.UpdateLinkRequest{
		ID:      req.ID,
		Changes: req.Changes,
	})
	if err != nil {
		errors.SendError(w, "failed to update link: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send response
	json.NewEncoder(w).Encode(responder.New(link))
}
