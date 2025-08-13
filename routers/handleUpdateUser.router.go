package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
	"github.com/hillview.tv/coreAPI/util"
)

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req query.UpdateUserRequest

	// get the query var
	q := mux.Vars(r)["query"]

	// parse the query var
	if q == "" {
		responder.ParamError(w, "query")
		return
	} else {
		intID, err := strconv.Atoi(q)
		if err != nil {
			responder.SendError(w, "query is not an ID", http.StatusBadRequest)
			return
		}

		req.ID = &intID
	}

	// parse the body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responder.SendError(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	// validate the fields
	if req.ID == nil {
		responder.ParamError(w, "id")
		return
	}

	if req.Changes == nil {
		responder.ParamError(w, "changes")
		return
	}

	// validate email formatting
	if req.Changes.Email != nil {
		if !util.IsValidEmail(*req.Changes.Email) {
			responder.SendError(w, "invalid email", http.StatusBadRequest)
			return
		}
	}

	// validate there is at least one change
	if req.Changes.Name == nil &&
		req.Changes.Email == nil &&
		req.Changes.Username == nil &&
		req.Changes.ProfileImageURL == nil &&
		req.Changes.Authentication == nil {
		responder.SendError(w, "no changes provided", http.StatusBadRequest)
		return
	}

	// run the query
	user, err := query.UpdateUser(db.DB, req)
	if err != nil {
		responder.SendError(w, "failed to update user: "+err.Error(), http.StatusConflict)
		return
	}

	// return the user
	responder.New(w, user)
}
