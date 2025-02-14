package routers

import (
	"encoding/json"
	"net/http"

	"github.com/hillview.tv/coreAPI/responder"
)

func HandleUpdateSpotlight(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(responder.Error("Not implemented"))
}
