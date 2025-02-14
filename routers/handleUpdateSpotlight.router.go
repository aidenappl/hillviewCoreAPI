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

func HandleUpdateSpotlight(w http.ResponseWriter, r *http.Request) {
	//  parse body
	req := query.UpdateSpotlightRequest{}

	// get rank from url
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
		req.Rank = &rank
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(responder.Error(err.Error()))
		return
	}

	//  update spotlight
	spotlight, err := query.UpdateSpotlight(db.DB, req)
	if err != nil {
		json.NewEncoder(w).Encode(responder.Error(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(responder.New(spotlight))
}
