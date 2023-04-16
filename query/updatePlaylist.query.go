package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type UpdatePlaylistRequest struct {
	ID      *int
	Changes *UpdatePlaylistChanges
}

type UpdatePlaylistChanges struct {
	Name        *string `json:"name"`
	Status      *int    `json:"status"`
	Description *string `json:"description"`
	BannerImage *string `json:"banner_image"`
	Route       *string `json:"route"`
}

func UpdatePlaylist(db db.Queryable, req UpdatePlaylistRequest) (*structs.Playlist, error) {
	// check required fields
	if req.ID == nil {
		return nil, fmt.Errorf("required field id is nil")
	}

	if req.Changes == nil {
		return nil, fmt.Errorf("required field changes is nil")
	}

	// build query
	q := sq.Update("playlists").
		Where(sq.Eq{"id": *req.ID})

	// add changes
	if req.Changes.Name != nil {
		q = q.Set("name", *req.Changes.Name)
	}

	if req.Changes.Status != nil {
		q = q.Set("status", req.Changes.Status)
	}

	if req.Changes.Description != nil {
		q = q.Set("description", *req.Changes.Description)
	}

	if req.Changes.BannerImage != nil {
		q = q.Set("banner_image", *req.Changes.BannerImage)
	}

	if req.Changes.Route != nil {
		q = q.Set("route", *req.Changes.Route)
	}

	// execute query
	_, err := q.RunWith(db).Exec()
	if err != nil {
		return nil, err
	}

	// get updated playlist
	return GetPlaylist(db, GetPlaylistRequest{ID: req.ID})
}
