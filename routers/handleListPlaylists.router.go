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

type HandleListPlaylistsRequest struct {
	Limit  *int
	Offset *int
	Search *string
	Sort   *string
}

func HandleListPlaylists(w http.ResponseWriter, r *http.Request) {
	var req HandleListPlaylistsRequest

	// get from query params
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")

	// parse query params
	if limit != "" {
		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			errors.ParamError(w, "limit")
			return
		}
		req.Limit = &intLimit
	} else {
		errors.ParamError(w, "limit")
		return
	}

	if offset != "" {
		intOffset, err := strconv.Atoi(offset)
		if err != nil {
			errors.ParamError(w, "offset")
			return
		}
		req.Offset = &intOffset
	} else {
		errors.ParamError(w, "offset")
		return
	}

	if search != "" {
		req.Search = &search
	}

	if sort != "" {
		req.Sort = &sort
	}

	// run the query
	playlists, err := query.ListPlaylists(db.DB, query.ListPlaylistsRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Search: req.Search,
		Sort:   req.Sort,
	})
	if err != nil {
		errors.SendError(w, "failed to list playlists: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send response
	json.NewEncoder(w).Encode(responder.New(playlists))
}
