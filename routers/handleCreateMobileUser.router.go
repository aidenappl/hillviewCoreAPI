package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
	"github.com/hillview.tv/coreAPI/util"
)

func HandleCreateMobileUser(w http.ResponseWriter, r *http.Request) {
	var req query.CreateMobileUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responder.SendError(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if req.Name == nil {
		responder.ParamError(w, "name")
		return
	}

	if req.Email == nil {
		responder.ParamError(w, "email")
		return
	} else {
		// validate email
		if !util.IsValidEmail(*req.Email) {
			responder.SendError(w, "invalid email", http.StatusBadRequest)
			return
		}
	}

	if req.Identifier == nil {
		responder.ParamError(w, "identifier")
		return
	}

	// run the query
	user, err := query.CreateMobileUser(db.AssetDB, req)
	if err != nil {
		responder.SendError(w, "failed to create mobile user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send the response
	responder.New(w, user)

}
