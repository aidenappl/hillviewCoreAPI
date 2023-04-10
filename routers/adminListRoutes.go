package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
)

type HandleListPlaylistsRequest struct {
	Limit  *uint64 `json:"limit"`
	Offset *uint64 `json:"offset"`
}

func HandleListPlaylists(w http.ResponseWriter, r *http.Request) {
	var req HandleListPlaylistsRequest

	// Get the Limit Param
	limit := r.URL.Query().Get("limit")
	if len(limit) == 0 {
		http.Error(w, "missing limit param", http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.ParseUint(string(limit), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Limit = &limitInt

	// Get the Offset Param
	offset := r.URL.Query().Get("offset")

	if len(offset) == 0 {
		http.Error(w, "missing offset param", http.StatusBadRequest)
		return
	}

	offsetInt, err := strconv.ParseUint(string(offset), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Offset = &offsetInt

	playlists, err := query.ListPlaylists(db.DB, query.ListPlaylistsRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(playlists)
}

func HandleListAdminUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	limit := params["limit"]

	if len(limit) == 0 {
		http.Error(w, "missing limit param", http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.ParseUint(string(limit), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	users, err := query.ListAdminUsers(db.DB, query.ListAdminUsersRequest{
		Limit: &limitInt,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func HandleListLinks(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")

	if len(limit) == 0 {
		http.Error(w, "missing limit param", http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.ParseUint(string(limit), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	links, err := query.ListLinks(db.DB, query.ListLinksRequest{
		Limit: &limitInt,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(links)

}

type HandleListMobileUsersRequest struct {
	Limit  *uint64 `json:"limit"`
	Offset *uint64 `json:"offset"`
	Search *string `json:"search"`
	Sort   *string `json:"sort"`
}

func HandleListMobileUsers(w http.ResponseWriter, r *http.Request) {

	var req HandleListMobileUsersRequest

	// Get the Sort param
	sort := r.URL.Query().Get("sort")
	if len(sort) != 0 {
		req.Sort = &sort
	}
	if req.Sort == nil {
		var forcedSort = "DESC"
		req.Sort = &forcedSort
	}

	// Get the Limit Param
	limit := r.URL.Query().Get("limit")
	if len(limit) != 0 {
		limitInt, err := strconv.ParseUint(string(limit), 10, 64)
		if err != nil {
			http.Error(w, "failed to convert string to int: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Limit = &limitInt
	}
	if req.Limit == nil {
		http.Error(w, "missing limit param", http.StatusBadRequest)
		return
	}

	// Get the Offset Param
	offset := r.URL.Query().Get("offset")
	if len(offset) != 0 {
		offsetInt, err := strconv.ParseUint(string(offset), 10, 64)
		if err != nil {
			http.Error(w, "failed to convert string to int: "+err.Error(), http.StatusInternalServerError)
			return
		}
		req.Offset = &offsetInt
	}
	if req.Offset == nil {
		http.Error(w, "missing offset param", http.StatusBadRequest)
		return
	}

	// Get the Search Param
	search := r.URL.Query().Get("search")
	if len(search) != 0 {
		req.Search = &search
	}

	users, err := query.GetMobileUsers(db.AssetDB, query.GetMobileUsersRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Search: req.Search,
		Sort:   req.Sort,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func HandleListCheckouts(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")

	if len(limit) == 0 {
		http.Error(w, "missing limit param", http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.ParseUint(string(limit), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	checkouts, err := query.ListCheckouts(db.AssetDB, query.ListCheckoutsRequest{
		Limit: &limitInt,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(checkouts)
}

func HandleListOpenCheckouts(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")

	if len(limit) == 0 {
		http.Error(w, "missing limit param", http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.ParseUint(string(limit), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	checkouts, err := query.ListOpenCheckouts(db.AssetDB, query.ListOpenCheckoutsRequest{
		Limit: &limitInt,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(checkouts)
}

func HandleListVideos(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	limit := params["limit"]

	if len(limit) == 0 {
		http.Error(w, "missing limit param", http.StatusBadRequest)
		return
	}

	limitInt, err := strconv.ParseUint(string(limit), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	videos, err := query.ListVideos(db.DB, query.ListVideosRequest{
		Limit:           limitInt,
		Offset:          0,
		IncludeArchived: true,
		IncludeDrafts:   true,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(videos)
}
