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

type HandleGetUserRequest struct {
	ID *int
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	var req HandleGetUserRequest

	// get the query var
	q := mux.Vars(r)["query"]

	// parse the query
	if q == "" {
		errors.ParamError(w, "query")
		return
	} else {
		queryID, err := strconv.Atoi(q)
		if err != nil {
			errors.SendError(w, "query must be an integer", http.StatusBadRequest)
			return
		}

		req.ID = &queryID
	}

	// run the query
	user, err := query.GetUser(db.DB, query.GetUserRequest{
		ID: req.ID,
	})
	if err != nil {
		errors.SendError(w, "failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		errors.SendError(w, "user not found", http.StatusNotFound)
		return
	}

	// send the response
	json.NewEncoder(w).Encode(responder.New(user))
}
