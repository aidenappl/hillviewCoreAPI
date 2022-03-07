package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
)

func HandleSearchAdminUsers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	if search == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := query.SearchMobileUsers(db.AssetDB, query.SearchMobileUsersRequest{
		Search: search,
		Limit:  nil,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}
