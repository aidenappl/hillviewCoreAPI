package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

func HandleReorderSpotlight(w http.ResponseWriter, r *http.Request) {
	var items []query.ReorderSpotlightItem

	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		responder.BadBody(w, err)
		return
	}

	if len(items) == 0 {
		responder.SendError(w, "request body must contain at least one item", http.StatusBadRequest)
		return
	}

	spotlights, err := query.ReorderSpotlights(db.DB, query.ReorderSpotlightsRequest{Items: items})
	if err != nil {
		responder.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responder.New(w, spotlights)
}
