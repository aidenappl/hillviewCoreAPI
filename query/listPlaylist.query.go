package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListPlaylistsRequest struct {
	Limit  *int
	Offset *int
	Search *string
	Sort   *string

	// select fields
	ID *int
}

func ListPlaylists(db db.Queryable, req ListPlaylistsRequest) ([]*structs.Playlist, error) {
	// check required fields
	if req.Limit == nil {
		return nil, fmt.Errorf("required field limit is nil")
	}

	if req.Offset == nil {
		return nil, fmt.Errorf("required field offset is nil")
	}

	// check sort formatting
	if req.Sort != nil {
		if *req.Sort != "desc" && *req.Sort != "asc" {
			return nil, fmt.Errorf("invalid sort provided")
		}
	}

	// build query
	q := sq.Select(
		"playlists.id",
		"playlists.name",
		"playlists.description",
		"playlists.banner_image",
		"playlists.route",
		"playlists.inserted_at",

		"playlist_statuses.id",
		"playlist_statuses.name",
		"playlist_statuses.short_name",
	).
		From("playlists").
		Join("playlist_statuses ON playlists.status = playlist_statuses.id").
		OrderBy("playlists.id DESC").
		Limit(uint64(*req.Limit)).
		Offset(uint64(*req.Offset))

	// add select fields
	if req.ID != nil {
		q = q.Where(sq.Eq{"playlists.id": *req.ID})
	}

	// add search
	if req.Search != nil {
		q = q.Where(sq.Like{"playlists.name": *req.Search})
	}

	// add sort
	if req.Sort != nil {
		q = q.OrderBy(fmt.Sprintf("playlists.id %s", *req.Sort))
	}

	// run query
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	// execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	// scan rows
	var playlists []*structs.Playlist
	for rows.Next() {
		var playlist structs.Playlist
		var status structs.GeneralNSM
		err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.Description,
			&playlist.BannerImage,
			&playlist.Route,
			&playlist.InsertedAt,

			&status.ID,
			&status.Name,
			&status.ShortName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Get the videos for the playlist
		videos, err := ListVideos(db, ListVideosRequest{
			Limit:      &[]int{100}[0],
			Offset:     &[]int{0}[0],
			PlaylistID: &playlist.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("error getting videos for playlist: %w", err)
		}
		playlist.Videos = videos

		playlist.Status = status

		playlists = append(playlists, &playlist)
	}

	return playlists, nil
}
