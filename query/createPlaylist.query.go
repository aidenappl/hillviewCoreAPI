package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	cdb "github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type CreatePlaylistRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	BannerImage *string `json:"banner_image"`
	Route       *string `json:"route"`
	Videos      *[]int  `json:"videos"`
}

func CreatePlaylist(db cdb.Queryable, req CreatePlaylistRequest) (*structs.Playlist, error) {
	// validate fields
	if req.Name == nil || *req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if req.Description == nil || *req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	if req.BannerImage == nil || *req.BannerImage == "" {
		return nil, fmt.Errorf("banner_image is required")
	}

	if req.Route == nil || *req.Route == "" {
		return nil, fmt.Errorf("route is required")
	}

	if req.Videos == nil || len(*req.Videos) == 0 {
		return nil, fmt.Errorf("videos is required")
	}

	// insert playlist
	query, args, err := sq.Insert("playlists").
		Columns(
			"name",
			"description",
			"banner_image",
			"route",
		).
		Values(
			*req.Name,
			*req.Description,
			*req.BannerImage,
			*req.Route,
		).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql playlist request: %w", err)
	}

	// execute query
	rows, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql playlist request: %w", err)
	}

	// get id
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// insert playlist videos
	for _, videoID := range *req.Videos {
		video, err := GetVideo(db, GetVideoRequest{
			ID: &videoID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get video: %w", err)
		}

		if video == nil {
			return nil, fmt.Errorf("video with id %d does not exist", videoID)
		}

		query, args, err := sq.Insert("playlist_associations").
			Columns(
				"playlist_id",
				"video_id",
			).
			Values(
				id,
				videoID,
			).
			ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build sql playlist video request: %w", err)
		}

		// execute query
		_, err = db.Exec(query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute sql playlist video request: %w", err)
		}
	}

	// convert id to int
	idInt := int(id)

	// get playlist
	return GetPlaylist(db, GetPlaylistRequest{
		ID: &idInt,
	})
}
