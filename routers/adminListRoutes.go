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

func HandleListMobileUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := query.ListMobileUsers(db.AssetDB, query.ListMobileUsersRequest{
		Limit: &limitInt,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func HandleListAssets(w http.ResponseWriter, r *http.Request) {
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

	assets, err := query.ListAssets(db.AssetDB, query.ListAssetsRequest{
		Limit: &limitInt,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(assets)
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
