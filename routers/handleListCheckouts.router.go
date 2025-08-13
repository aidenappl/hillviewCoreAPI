package routers

import (
	"net/http"
	"strconv"

	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleListCheckoutsRequest struct {
	Limit  int
	Offset int
	Sort   *string
}

func HandleListCheckouts(w http.ResponseWriter, r *http.Request) {
	var req HandleListCheckoutsRequest

	// get from query params
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	sort := r.URL.Query().Get("sort")

	// parse query params
	if limit != "" {
		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			responder.SendError(w, "invalid limit provided", http.StatusBadRequest)
			return
		}
		req.Limit = intLimit
	} else {
		responder.ErrRequiredKey(w, "limit")
		return
	}

	if offset != "" {
		intOffset, err := strconv.Atoi(offset)
		if err != nil {
			responder.SendError(w, "invalid offset provided", http.StatusBadRequest)
			return
		}
		req.Offset = intOffset
	} else {
		responder.ErrRequiredKey(w, "offset")
		return
	}

	if sort != "" {
		req.Sort = &sort
	}

	// get checkouts
	checkouts, err := query.ListCheckouts(db.AssetDB, query.ListCheckoutsRequest{
		Limit:  &req.Limit,
		Offset: &req.Offset,
		Sort:   req.Sort,
	})
	if err != nil {
		responder.SendError(w, "failed to list checkouts: "+err.Error(), http.StatusConflict)
		return
	}

	// respond
	responder.New(w, checkouts)
}
