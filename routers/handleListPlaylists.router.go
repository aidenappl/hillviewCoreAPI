package routers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleListPlaylistsRequest struct {
	Limit    *int
	Offset   *int
	Search   *string
	Sort     *string
	SortBy   *string
	Statuses *[]int
}

func HandleListPlaylists(w http.ResponseWriter, r *http.Request) {
	var req HandleListPlaylistsRequest

	// get from query params
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")
	sortBy := r.URL.Query().Get("sort_by")
	status := r.URL.Query().Get("status")

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

	// Parse comma-separated status IDs, e.g. "1,3"
	if status != "" {
		parts := strings.Split(status, ",")
		statuses := make([]int, 0, len(parts))
		for _, p := range parts {
			id, err := strconv.Atoi(strings.TrimSpace(p))
			if err == nil {
				statuses = append(statuses, id)
			}
		}
		if len(statuses) > 0 {
			req.Statuses = &statuses
		}
	}

	// run the query
	playlists, err := query.ListPlaylists(db.DB, query.ListPlaylistsRequest{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Search:   req.Search,
		Sort:     req.Sort,
		SortBy:   req.SortBy,
		Statuses: req.Statuses,
	})
	if err != nil {
		responder.SendError(w, "failed to list playlists: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send response
	responder.New(w, playlists)
}
