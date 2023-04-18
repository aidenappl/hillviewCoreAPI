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

type HandleGetMobileUserRequest struct {
	ID         *int
	Identifier *string
}

func HandleGetMobileUser(w http.ResponseWriter, r *http.Request) {
	var req HandleGetMobileUserRequest

	// parse the query params
	q := mux.Vars(r)["query"]

	// set the query params
	if q == "" {
		errors.ParamError(w, "query")
		return
	} else {
		idInt, err := strconv.Atoi(q)
		if err != nil {
			req.Identifier = &q
		} else {
			req.ID = &idInt
		}
	}

	// run the query
	mobileUser, err := query.GetMobileUser(db.AssetDB, query.GetMobileUserRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
	})
	if err != nil {
		errors.SendError(w, "failed to get mobile user: "+err.Error(), http.StatusConflict)
		return
	}

	if mobileUser == nil {
		errors.SendError(w, "mobile user not found", http.StatusNotFound)
		return
	}

	// send the response
	json.NewEncoder(w).Encode(responder.New(mobileUser))
}
