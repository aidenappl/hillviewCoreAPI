package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
)

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

	checkouts, err := query.ListCheckouts(db.AssetDB, query.ListCheckoutsRequest{
		Limit: &limitInt,
	})
	if err != nil {
		http.Error(w, "failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(checkouts)
}
