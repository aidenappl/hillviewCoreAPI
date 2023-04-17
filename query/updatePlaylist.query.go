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
	Name         *string `json:"name"`
	Status       *int    `json:"status"`
	Description  *string `json:"description"`
	BannerImage  *string `json:"banner_image"`
	Route        *string `json:"route"`
	RemoveVideos *[]int  `json:"remove_videos"`
	AddVideos    *[]int  `json:"add_videos"`
}

func UpdatePlaylist(db db.Queryable, req UpdatePlaylistRequest) (*structs.Playlist, error) {
	changesMade := false

	// check required fields
	if req.ID == nil {
		return nil, fmt.Errorf("required field id is nil")
	}

	if req.Changes == nil {
		return nil, fmt.Errorf("required field changes is nil")
	}

	if req.Changes.RemoveVideos != nil {
		changesMade = true
		// remove videos
		for _, videoID := range *req.Changes.RemoveVideos {
			query, args, err := sq.Delete("playlist_associations").
				Where(sq.Eq{"playlist_id": *req.ID, "video_id": videoID}).
				ToSql()
			if err != nil {
				return nil, fmt.Errorf("failed to form query: %v", err)
			}

			_, err = db.Exec(query, args...)
			if err != nil {
				return nil, fmt.Errorf("failed to remove video from playlist: %v", err)
			}
		}
	}

	if req.Changes.AddVideos != nil {
		changesMade = true
		// add videos
		for _, videoID := range *req.Changes.AddVideos {
			query, args, err := sq.Insert("playlist_associations").
				Columns("playlist_id", "video_id").
				Values(*req.ID, videoID).
				ToSql()
			if err != nil {
				return nil, fmt.Errorf("failed to form query: %v", err)
			}

			_, err = db.Exec(query, args...)
			if err != nil {
				return nil, fmt.Errorf("failed to add video to playlist: %v", err)
			}
		}
	}

	if changesMade && req.Changes.Name == nil && req.Changes.Status == nil && req.Changes.Description == nil && req.Changes.BannerImage == nil && req.Changes.Route == nil {
		// if no changes were made to the playlist, just return the playlist
		return GetPlaylist(db, GetPlaylistRequest{ID: req.ID})
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

	// form query
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to form query: %v", err)
	}

	// execute query
	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// get updated playlist
	return GetPlaylist(db, GetPlaylistRequest{ID: req.ID})
}
