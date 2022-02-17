package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/structs"
)

type DeleteVideoRequest struct {
	ID *int `json:"id"`
}

func HandleDeleteVideo(w http.ResponseWriter, r *http.Request) {
	body := DeleteVideoRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.ID == nil {
		http.Error(w, "no video id provided", http.StatusBadRequest)
		return
	}

	_, err = query.EditVideo(db.DB, query.EditVideoRequest{
		ID: body.ID,
		Modifications: &query.VideoModifications{
			Status: &structs.VideoStatusDeleted,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
