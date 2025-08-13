package routers

import (
	"net/http"
	"strconv"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

func HandleListSpotlight(w http.ResponseWriter, r *http.Request) {

	q := query.ListSpotlightsRequest{}

	// get params
	params := r.URL.Query()
	limit := params.Get("limit")
	if limit == "" {
		responder.ParamError(w, "limit")
		return
	} else {
		l, err := strconv.Atoi(limit)
		if err != nil {
			responder.SendError(w, "invalid limit", http.StatusBadRequest)
			return
		}
		q.Limit = &l
	}
	offset := params.Get("offset")
	if offset == "" {
		responder.ParamError(w, "offset")
		return
	} else {
		o, err := strconv.Atoi(offset)
		if err != nil {
			responder.SendError(w, "invalid offset", http.StatusBadRequest)
			return
		}
		q.Offset = &o
	}

	spotlights, err := query.ListSpotlights(db.DB, q)
	if err != nil {
		responder.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responder.New(w, spotlights)
}
