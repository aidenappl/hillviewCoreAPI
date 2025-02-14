package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

func HandleGetSpotlight(w http.ResponseWriter, r *http.Request) {
	q := query.GetSpotlightRequest{}

	// get the query var
	rankVar := mux.Vars(r)["rank"]
	if rankVar == "" {
		json.NewEncoder(w).Encode(responder.Error("missing rank"))
		return
	} else {
		rank, err := strconv.Atoi(rankVar)
		if err != nil {
			json.NewEncoder(w).Encode(responder.Error("invalid rank"))
			return
		}
		q.Rank = &rank
	}

	spotlight, err := query.GetSpotlight(db.DB, q)
	if err != nil {
		json.NewEncoder(w).Encode(responder.Error(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(responder.New(spotlight))
}
