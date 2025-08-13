package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/middleware"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleCreateLinkRequest struct {
	query.CreateLinkRequest
}

func HandleCreateLink(w http.ResponseWriter, r *http.Request) {
	var req HandleCreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder.SendError(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// validate the request

	// add user id
	user := middleware.WithUserModelValue(r.Context())
	req.CreateLinkRequest.CreatedBy = user.ID

	// create the link
	link, err := query.CreateLink(db.DB, req.CreateLinkRequest)
	if err != nil {
		responder.SendError(w, "failed to create link: "+err.Error(), http.StatusConflict)
		return
	}

	// send the response
	responder.New(w, link)
}
