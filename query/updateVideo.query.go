package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type UpdateVideoRequest struct {
	// Ident Fields from Query
	ID         *int    `json:"id"`
	Identifier *string `json:"identifier"`

	// Edit Fields
	Changes *UpdateVideoChanges `json:"changes"`
}

type UpdateVideoChanges struct {
	// Asset Fields
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	Thumbnail      *string `json:"thumbnail"`
	AllowDownloads *bool   `json:"allow_downloads"`
	DownloadURL    *string `json:"download_url"`
	URL            *string `json:"url"`
	Status         *int    `json:"status"`
}

func UpdateVideo(db db.Queryable, req UpdateVideoRequest) (*structs.Video, error) {
	// check required fields
	if req.ID == nil && req.Identifier == nil {
		return nil, fmt.Errorf("id or identifier is required")
	}

	// check changes
	if req.Changes == nil {
		return nil, fmt.Errorf("changes is required")
	}

	// check that there is at least one change
	if req.Changes.Title == nil && req.Changes.Description == nil && req.Changes.Thumbnail == nil && req.Changes.AllowDownloads == nil && req.Changes.DownloadURL == nil && req.Changes.URL == nil && req.Changes.Status == nil {
		return nil, fmt.Errorf("no changes provided")
	}

	q := sq.Update("videos")

	// set fields
	if req.Changes.Title != nil {
		q = q.Set("title", req.Changes.Title)
	}

	if req.Changes.Description != nil {
		q = q.Set("description", req.Changes.Description)
	}

	if req.Changes.Thumbnail != nil {
		q = q.Set("thumbnail", req.Changes.Thumbnail)
	}

	if req.Changes.AllowDownloads != nil {
		q = q.Set("allow_downloads", req.Changes.AllowDownloads)
	}

	if req.Changes.DownloadURL != nil {
		q = q.Set("download_url", req.Changes.DownloadURL)
	}

	if req.Changes.URL != nil {
		q = q.Set("url", req.Changes.URL)
	}

	if req.Changes.Status != nil {
		q = q.Set("status", req.Changes.Status)
	}

	// set where
	if req.ID != nil {
		q = q.Where(sq.Eq{"id": req.ID})
	}

	if req.Identifier != nil {
		q = q.Where(sq.Eq{"identifier": req.Identifier})
	}

	// build query
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	// run query
	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %v", err)
	}

	return GetVideo(db, GetVideoRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
	})
}
