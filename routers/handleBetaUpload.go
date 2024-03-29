package routers

import (
	"net/http"

	"github.com/hillview.tv/coreAPI/env"
)

func HandleBetaUpload(w http.ResponseWriter, r *http.Request) {
	cloudflareEndpoint := "https://api.cloudflare.com/client/v4/accounts/" + env.CloudflareAccountID + "/stream?direct_user=true"
	client := &http.Client{}

	req, err := http.NewRequest("POST", cloudflareEndpoint, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", "bearer "+env.CloudflareAuthToken)
	req.Header.Add("Tus-Resumable", "1.0.0")
	req.Header.Add("Upload-Length", r.Header.Get("Upload-Length"))
	req.Header.Add("Upload-Metadata", r.Header.Get("Upload-Metadata"))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusSeeOther)
}
