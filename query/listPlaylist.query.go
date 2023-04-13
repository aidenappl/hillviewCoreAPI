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
	).
		From("playlists").
		OrderBy("playlists.id DESC").
		Limit(uint64(*req.Limit)).
		Offset(uint64(*req.Offset))

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
		err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.Description,
			&playlist.BannerImage,
			&playlist.Route,
			&playlist.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		playlists = append(playlists, &playlist)
	}

	return playlists, nil
}
