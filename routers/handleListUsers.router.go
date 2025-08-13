package routers

import (
	"net/http"
	"strconv"

	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleListUsersRequest struct {
	Limit  *int
	Offset *int
	Search *string
	Sort   *string
}

func HandleListUsers(w http.ResponseWriter, r *http.Request) {
	req := HandleListUsersRequest{}

	// parse the query params
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")

	// set the query params
	if limit != "" {
		limitID, err := strconv.Atoi(limit)
		if err != nil {
			responder.SendError(w, "limit must be an integer", http.StatusBadRequest)
			return
		}

		req.Limit = &limitID
	} else {
		responder.ParamError(w, "limit")
		return
	}

	if offset != "" {
		offsetID, err := strconv.Atoi(offset)
		if err != nil {
			responder.SendError(w, "offset must be an integer", http.StatusBadRequest)
			return
		}

		req.Offset = &offsetID
	} else {
		responder.ParamError(w, "offset")
		return
	}

	if search != "" {
		req.Search = &search
	}

	if sort != "" {
		req.Sort = &sort
	}

	// run the query
	users, err := query.ListUsers(db.DB, query.ListUsersRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Search: req.Search,
		Sort:   req.Sort,
	})
	if err != nil {
		responder.SendError(w, "failed to list users: "+err.Error(), http.StatusConflict)
		return
	}

	// send the response
	responder.New(w, users)
}
