package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/structs"
)

type CreateLinkRequest struct {
	Route    *string `json:"route"`
	Endpoint *string `json:"endpoint"`
}

func CreateLink(w http.ResponseWriter, r *http.Request) {
	body := CreateLinkRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if body.Route == nil || len(*body.Route) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if body.Endpoint == nil || len(*body.Endpoint) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := r.Context().Value("user").(structs.User)

	err = query.CreateLink(db.DB, query.CreateLinkRequest{
		Route:    body.Route,
		Endpoint: body.Endpoint,
		Creator:  &user.ID,
	})
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
