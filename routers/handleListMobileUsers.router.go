package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/errors"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleListMobileUsersRequest struct {
	Limit  *int
	Offset *int
	Search *string
	Sort   *string
}

func HandleListMobileUsers(w http.ResponseWriter, r *http.Request) {
	var req HandleListMobileUsersRequest

	// parse the query params
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")

	// set the query params
	if limit == "" {
		errors.ParamError(w, "limit")
		return
	} else {
		limitID, err := strconv.Atoi(limit)
		if err != nil {
			errors.SendError(w, "limit must be an integer", http.StatusBadRequest)
			return
		}

		req.Limit = &limitID
	}

	if offset == "" {
		errors.ParamError(w, "offset")
		return
	} else {
		offsetID, err := strconv.Atoi(offset)
		if err != nil {
			errors.SendError(w, "offset must be an integer", http.StatusBadRequest)
			return
		}

		req.Offset = &offsetID
	}

	if search != "" {
		req.Search = &search
	}

	if sort != "" {
		req.Sort = &sort
	}

	// run the query
	mobileUsers, err := query.ListMobileUsers(db.AssetDB, query.ListMobileUsersRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Search: req.Search,
		Sort:   req.Sort,
	})
	if err != nil {
		errors.SendError(w, "failed to list mobile users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send the response
	json.NewEncoder(w).Encode(responder.New(mobileUsers))
}
