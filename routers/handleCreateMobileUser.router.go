package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/errors"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

func HandleCreateMobileUser(w http.ResponseWriter, r *http.Request) {
	var req query.CreateMobileUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.SendError(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	// TODO: Validate the request

	// run the query
	user, err := query.CreateMobileUser(db.AssetDB, req)
	if err != nil {
		errors.SendError(w, "failed to create mobile user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send the response
	json.NewEncoder(w).Encode(responder.New(user))

}
