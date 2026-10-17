package routers

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hillview.tv/coreAPI/responder"
	"github.com/hillview.tv/coreAPI/util"
)

type HandleImageUploadResponse struct {
	URL string `json:"url"`
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		responder.SendError(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest, err)
		return
	}

	// Check ID param
	var idInt int
	id := r.FormValue("id")
	if id == "" {
		responder.SendError(w, "id is required", http.StatusBadRequest)
		return
	}

	idInt, err = strconv.Atoi(id)
	if err != nil {
		responder.SendError(w, "Invalid id parameter", http.StatusBadRequest, err)
		return
	}

	// Check route param
	route := r.FormValue("route")
	if route == "" {
		responder.SendError(w, "route is required", http.StatusBadRequest)
		return
	}

	validRoute := util.ValidS3Route(route)
	if !validRoute {
		responder.SendError(w, "Invalid route, must be in the permitted routes case", http.StatusBadRequest)
		return
	}

	// Get the image
	file, _, err := r.FormFile("image")
	if err != nil {
		responder.SendError(w, "image is required", http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		responder.SendError(w, "failed to read image", http.StatusConflict, err)
		return
	}

	url, err := util.UploadImage(context.Background(), util.UploadImageRequest{
		Image: fileBytes,
		ID:    idInt,
		Route: route,
	})
	if err != nil {
		responder.SendError(w, "failed to upload image to s3", http.StatusConflict, err)
		return
	}

	responder.New(w, HandleImageUploadResponse{
		URL: url,
	})
}
