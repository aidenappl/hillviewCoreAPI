package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type EditVideoRequest struct {
	ID            *int                `json:"id"`
	Modifications *VideoModifications `json:"modifications"`
}

type VideoModifications struct {
	Title       *string `json:"title"`
	Thumbnail   *string `json:"thumbnail"`
	URL         *string `json:"url"`
	Description *string `json:"description"`
	Status      *int    `json:"status"`
}

func EditVideo(db db.Queryable, req EditVideoRequest) (*structs.Video, error) {
	if req.Modifications == nil {
		return nil, fmt.Errorf("no modifications provided")
	}

	if req.ID == nil {
		return nil, fmt.Errorf("no asset id provided")
	}

	dataToSet := map[string]interface{}{}

	if req.Modifications.Description != nil {
		dataToSet["description"] = *req.Modifications.Description
	}

	if req.Modifications.Thumbnail != nil {
		dataToSet["thumbnail"] = *req.Modifications.Thumbnail
	}

	if req.Modifications.Title != nil {
		dataToSet["title"] = *req.Modifications.Title
	}

	if req.Modifications.URL != nil {
		dataToSet["url"] = *req.Modifications.URL
	}

	if req.Modifications.Status != nil {
		dataToSet["status"] = *req.Modifications.Status
	}

	query, args, err := sq.Update("videos").
		SetMap(dataToSet).
		Where(sq.Eq{"id": *req.ID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}

	video, err := GetVideo(db, GetVideoRequest{
		ID: req.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	return video, nil
}
