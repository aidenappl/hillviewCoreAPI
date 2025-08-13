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
		responder.ParamError(w, "rank")
		return
	} else {
		rank, err := strconv.Atoi(rankVar)
		if err != nil {
			responder.SendError(w, "invalid rank", http.StatusBadRequest)
			return
		}
		req.Rank = &rank
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responder.BadBody(w, err)
		return
	}

	//  update spotlight
	spotlight, err := query.UpdateSpotlight(db.DB, req)
	if err != nil {
		responder.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responder.New(w, spotlight)
}
