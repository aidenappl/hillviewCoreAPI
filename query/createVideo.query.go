package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type CreateVideoRequest struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	Thumbnail      *string `json:"thumbnail"`
	URL            *string `json:"url"`
	DownloadURL    *string `json:"download_url"`
	AllowDownloads *bool   `json:"allow_downloads"`
	Status         *int    `json:"status"`
}

func CreateVideo(db db.Queryable, req CreateVideoRequest) (*structs.Video, error) {
	// validate required fields
	if req.Title == nil || *req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if req.Description == nil || *req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	if req.Thumbnail == nil || *req.Thumbnail == "" {
		return nil, fmt.Errorf("thumbnail is required")
	}

	if req.URL == nil || *req.URL == "" {
		return nil, fmt.Errorf("url is required")
	}

	if req.DownloadURL == nil || *req.DownloadURL == "" {
		return nil, fmt.Errorf("download_url is required")
	}

	if req.AllowDownloads == nil {
		return nil, fmt.Errorf("allow_downloads is required")
	}

	if req.Status == nil {
		return nil, fmt.Errorf("status is required")
	}

	// create the video
	query, args, err := sq.Insert("videos").
		Columns(
			"title",
			"description",
			"thumbnail",
			"url",
			"download_url",
			"allow_downloads",
			"status").
		Values(
			*req.Title,
			*req.Description,
			*req.Thumbnail,
			*req.URL,
			*req.DownloadURL,
			*req.AllowDownloads,
			*req.Status).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create sql query: %w", err)
	}

	// execute the query
	rows, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}

	// get the id of the newly created video
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// convert the id
	videoID := int(id)

	return GetVideo(db, GetVideoRequest{
		ID: &videoID,
	})
}
