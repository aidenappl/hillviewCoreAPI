package routers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hillview.tv/coreAPI/errors"
	"github.com/hillview.tv/coreAPI/responder"
	"github.com/hillview.tv/coreAPI/util"
)

type HandleImageUploadResponse struct {
	URL string `json:"url"`
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		errors.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check ID param
	var idInt int
	id := r.FormValue("id")
	if id == "" {
		errors.SendError(w, "id is required", http.StatusBadRequest)
		return
	}

	idInt, err = strconv.Atoi(id)
	if err != nil {
		errors.SendError(w, "Invalid id parameter: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check route param
	route := r.FormValue("route")
	if route == "" {
		errors.SendError(w, "route is required", http.StatusBadRequest)
		return
	}

	validRoute := util.ValidS3Route(route)
	if !validRoute {
		errors.SendError(w, "Invalid route, must be in the permitted routes case", http.StatusBadRequest)
		return
	}

	// Get the image
	file, _, err := r.FormFile("image")
	if err != nil {
		errors.SendError(w, "image is required: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		errors.SendError(w, "failed to read image: "+err.Error(), http.StatusConflict)
		return
	}

	url, err := util.UploadImage(context.Background(), util.UploadImageRequest{
		Image: fileBytes,
		ID:    idInt,
		Route: route,
	})
	if err != nil {
		errors.SendError(w, "failed to upload image to s3: "+err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(responder.New(HandleImageUploadResponse{
		URL: url,
	}))
}
