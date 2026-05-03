package routers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hillview.tv/coreAPI/db"

	"github.com/hillview.tv/coreAPI/query"
	"github.com/hillview.tv/coreAPI/responder"
)

type HandleListVideoRequest struct {
	Limit         int
	Offset        int
	Search        *string
	Sort          *string
	SortBy        *string
	Statuses      *[]int
	CreatorUserID *int
}

func HandleListVideo(w http.ResponseWriter, r *http.Request) {
	var req HandleListVideoRequest

	// get from query params
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")
	sortBy := r.URL.Query().Get("sort_by")
	status := r.URL.Query().Get("status")

	// parse query params
	if limit != "" {
		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			responder.SendError(w, "invalid limit provided", http.StatusBadRequest)
			return
		}
		req.Limit = intLimit
	} else {
		responder.ErrRequiredKey(w, "limit")
		return
	}

	if offset != "" {
		intOffset, err := strconv.Atoi(offset)
		if err != nil {
			responder.SendError(w, "invalid offset provided", http.StatusBadRequest)
			return
		}
		req.Offset = intOffset
	} else {
		responder.ErrRequiredKey(w, "offset")
		return
	}

	if search != "" {
		req.Search = &search
	}

	if sort != "" {
		req.Sort = &sort
	}

	if sortBy != "" {
		req.SortBy = &sortBy
	}

	// Parse comma-separated status IDs, e.g. "1,2,3"
	if status != "" {
		parts := strings.Split(status, ",")
		statuses := make([]int, 0, len(parts))
		for _, p := range parts {
			id, err := strconv.Atoi(strings.TrimSpace(p))
			if err == nil {
				statuses = append(statuses, id)
			}
		}
		if len(statuses) > 0 {
			req.Statuses = &statuses
		}
	}

	userID := r.URL.Query().Get("user_id")
	if userID != "" {
		intUserID, err := strconv.Atoi(userID)
		if err == nil {
			req.CreatorUserID = &intUserID
		}
	}

	// get the list of videos
	videos, err := query.ListVideos(db.DB, query.ListVideosRequest{
		Limit:         &req.Limit,
		Offset:        &req.Offset,
		Search:        req.Search,
		Sort:          req.Sort,
		SortBy:        req.SortBy,
		Statuses:      req.Statuses,
		CreatorUserID: req.CreatorUserID,
	})
	if err != nil {
		responder.SendError(w, "failed to list videos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// send the response
	responder.New(w, videos)
}
