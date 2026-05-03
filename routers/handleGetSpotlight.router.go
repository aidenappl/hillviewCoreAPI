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
	positionVar := mux.Vars(r)["position"]
	if positionVar == "" {
		responder.ParamError(w, "missing position")
		return
	} else {
		position, err := strconv.Atoi(positionVar)
		if err != nil {
			responder.ParamError(w, "invalid position")
			return
		}
		q.Position = &position
	}

	spotlight, err := query.GetSpotlight(db.DB, q)
	if err != nil {
		responder.ParamError(w, err.Error())
		return
	}

	responder.New(w, spotlight)
}
