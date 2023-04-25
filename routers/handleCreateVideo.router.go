package routers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/errors"
	"github.com/hillview.tv/coreAPI/mailer"
	"github.com/hillview.tv/coreAPI/middleware"
	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleCreateVideoRequest struct {
	query.CreateVideoRequest
}

func HandleCreateVideo(w http.ResponseWriter, r *http.Request) {
	// get the user from context
	user := middleware.WithUserModelValue(r.Context())
	if user == nil {
		errors.SendError(w, "failed to get user from context", http.StatusInternalServerError)
		return
	}

	var req HandleCreateVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.SendError(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// validate the request
	if req.Title == nil || *req.Title == "" {
		errors.ParamError(w, "title")
		return
	}

	if req.Description == nil || *req.Description == "" {
		errors.ParamError(w, "description")
		return
	}

	if req.Thumbnail == nil || *req.Thumbnail == "" {
		errors.ParamError(w, "thumbnail")
		return
	}

	if req.URL == nil || *req.URL == "" {
		errors.ParamError(w, "url")
		return
	}

	// create the video
	video, err := query.CreateVideo(db.DB, req.CreateVideoRequest)
	if err != nil {
		errors.SendError(w, "failed to create video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// check if user is student. If so, send notification to admin
	if user.Authentication.ShortName == "student" {
		adminUsers, err := query.ListUsers(db.DB, query.ListUsersRequest{
			Limit:          &[]int{100}[0],
			Offset:         &[]int{0}[0],
			Authentication: &[]int{3}[0],
		})
		if err != nil {
			log.Println("failed to list admin users: " + err.Error())
			return
		}

		// send email to each admin
		for _, admin := range adminUsers {
			log.Println("ℹ️ sending email to admin", admin.Email)
			_, err = mailer.SendTemplate(mailer.SendTemplateRequest{
				ToEmail: admin.Email,
				ToName:  admin.Email,

				FromEmail: "notifications@hillview.tv",
				FromName:  "HillviewTV Notifications",

				TemplateID: "d-75bbe45074674af6b480110344af7091",
				DynamicData: map[string]interface{}{
					"title":             "New Upload Alert",
					"full_name":         admin.Name,
					"body":              "This email is to notify you that <b>" + user.Name + "</b> uploaded a video to the Hillview TV website.\n\nIt is currently in the drafts pending your approval. Please see the team dashboard for more information.",
					"action_button_url": "https://team.hillview.tv/team/dashboard/videos?inspect=" + video.UUID,
				},
			})
			if err != nil {
				log.Println("failed to send email: " + err.Error())
			}
		}
	}

	// send the response
	json.NewEncoder(w).Encode(responder.New(video))

}
