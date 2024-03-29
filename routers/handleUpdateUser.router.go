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
	"github.com/hillview.tv/coreAPI/util"
)

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req query.UpdateUserRequest

	// get the query var
	q := mux.Vars(r)["query"]

	// parse the query var
	if q == "" {
		errors.ParamError(w, "query")
		return
	} else {
		intID, err := strconv.Atoi(q)
		if err != nil {
			errors.SendError(w, "query is not an ID", http.StatusBadRequest)
			return
		}

		req.ID = &intID
	}

	// parse the body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.SendError(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	// validate the fields
	if req.ID == nil {
		errors.ParamError(w, "id")
		return
	}

	if req.Changes == nil {
		errors.ParamError(w, "changes")
		return
	}

	// validate email formatting
	if req.Changes.Email != nil {
		if !util.IsValidEmail(*req.Changes.Email) {
			errors.SendError(w, "invalid email", http.StatusBadRequest)
			return
		}
	}

	// validate there is at least one change
	if req.Changes.Name == nil &&
		req.Changes.Email == nil &&
		req.Changes.Username == nil &&
		req.Changes.ProfileImageURL == nil &&
		req.Changes.Authentication == nil {
		errors.SendError(w, "no changes provided", http.StatusBadRequest)
		return
	}

	// run the query
	user, err := query.UpdateUser(db.DB, req)
	if err != nil {
		errors.SendError(w, "failed to update user: "+err.Error(), http.StatusConflict)
		return
	}

	// return the user
	json.NewEncoder(w).Encode(responder.New(user))
}
