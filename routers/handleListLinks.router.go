package routers

import (
	"net/http"
	"strconv"

	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleListLinksRequest struct {
	Limit        *int
	Offset       *int
	Search       *string
	Sort         *string
	SortBy       *string
	ShowArchived *bool
}

func HandleListLinks(w http.ResponseWriter, r *http.Request) {
	var req HandleListLinksRequest

	// get from query params
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")
	sortBy := r.URL.Query().Get("sort_by")
	showArchived := r.URL.Query().Get("show_archived")

	// parse query params
	if limit != "" {
		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			responder.ParamError(w, "limit")
			return
		}
		req.Limit = &intLimit
	} else {
		responder.ParamError(w, "limit")
		return
	}

	if offset != "" {
		intOffset, err := strconv.Atoi(offset)
		if err != nil {
			responder.ParamError(w, "offset")
			return
		}
		req.Offset = &intOffset
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

	if sortBy != "" {
		req.SortBy = &sortBy
	}

	if showArchived == "true" {
		t := true
		req.ShowArchived = &t
	}

	// run the query
	links, err := query.ListLinks(db.DB, query.ListLinksRequest{
		Limit:        req.Limit,
		Offset:       req.Offset,
		Search:       req.Search,
		Sort:         req.Sort,
		SortBy:       req.SortBy,
		ShowArchived: req.ShowArchived,
	})
	if err != nil {
		responder.SendError(w, "failed to list links: "+err.Error(), http.StatusConflict)
		return
	}

	// send response
	responder.New(w, links)

}
