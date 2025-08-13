package routers

import (
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
		responder.ParamError(w, "missing rank")
		return
	} else {
		rank, err := strconv.Atoi(rankVar)
		if err != nil {
			responder.ParamError(w, "invalid rank")
			return
		}
		q.Rank = &rank
	}

	spotlight, err := query.GetSpotlight(db.DB, q)
	if err != nil {
		responder.ParamError(w, err.Error())
		return
	}

	responder.New(w, spotlight)
}
