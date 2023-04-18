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

func HandleUpdateCheckout(w http.ResponseWriter, r *http.Request) {
	var req query.UpdateCheckoutRequest

	// get the query var
	q := mux.Vars(r)["query"]

	// parse the query var
	if q == "" {
		errors.ParamError(w, "query")
		return
	} else {
		idQuery, err := strconv.Atoi(q)
		if err != nil {
			errors.ParamError(w, "query")
			return
		}
		req.ID = &idQuery
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.SendError(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if req.ID == nil {
		errors.ParamError(w, "id")
		return
	}

	if req.Changes == nil {
		errors.ParamError(w, "changes")
		return
	}

	// run the query
	checkout, err := query.UpdateCheckout(db.AssetDB, req)
	if err != nil {
		errors.SendError(w, "failed to update checkout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(responder.New(checkout))
}
