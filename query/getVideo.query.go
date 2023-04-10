package query

import (
	"fmt"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type GetVideoRequest struct {
	ID         *int
	Identifier *string
}

func GetVideo(db db.Queryable, req GetVideoRequest) (*structs.Video, error) {
	// validate the request
	if req.ID == nil && req.Identifier == nil {
		return nil, fmt.Errorf("no id or identifier provided")
	}

	// run list query to get the video
	videos, err := ListVideos(db, ListVideosRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
		UseOr:      true,
		Limit:      &[]int{1}[0],
		Offset:     &[]int{0}[0],
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// check if the video exists
	if len(videos) == 0 {
		return nil, nil
	}

	return videos[0], nil
}
