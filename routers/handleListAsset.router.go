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

type HandleListAssetRequest struct {
	Limit  int
	Offset int
	Search *string
	Sort   *string
}

func HandleListAsset(w http.ResponseWriter, r *http.Request) {
	var req HandleListAssetRequest

	// get from query params
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")

	// parse query params
	if limit != "" {
		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			errors.SendError(w, "invalid limit provided", http.StatusBadRequest)
			return
		}
		req.Limit = intLimit
	} else {
		errors.ErrRequiredKey(w, "limit")
		return
	}

	if offset != "" {
		intOffset, err := strconv.Atoi(offset)
		if err != nil {
			errors.SendError(w, "invalid offset provided", http.StatusBadRequest)
			return
		}
		req.Offset = intOffset
	} else {
		errors.ErrRequiredKey(w, "offset")
		return
	}

	if search != "" {
		req.Search = &search
	}

	if sort != "" {
		req.Sort = &sort
	}

	assets, err := query.ListAssets(db.AssetDB, query.ListAssetsRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Search: req.Search,
		Sort:   req.Sort,
	})
	if err != nil {
		errors.SendError(w, "failed to list assets: "+err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(responder.New(assets))
}
