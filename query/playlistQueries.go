package query

// type ListPlaylistsRequest struct {
// 	Limit  *uint64 `json:"limit"`
// 	Offset *uint64 `json:"offset"`
// }

// func ListPlaylists(db db.Queryable, req ListPlaylistsRequest) (*[]structs.Playlist, error) {
// 	q := sq.Select(
// 		"playlists.id",
// 		"playlists.name",
// 		"playlists.description",
// 		"playlists.banner_image",
// 		"playlists.route",
// 		"playlists.inserted_at",
// 	).
// 		From("playlists").
// 		OrderBy("playlists.id DESC").
// 		Limit(*req.Limit).
// 		Offset(*req.Offset)

// 	query, args, err := q.ToSql()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to build query: %v", err)
// 	}

// 	rows, err := db.Query(query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute query: %v", err)
// 	}

// 	defer rows.Close()
// 	playlists := []structs.Playlist{}

// 	for rows.Next() {
// 		var playlist structs.Playlist
// 		err := rows.Scan(
// 			&playlist.ID,
// 			&playlist.Name,
// 			&playlist.Description,
// 			&playlist.BannerImage,
// 			&playlist.Route,
// 			&playlist.InsertedAt,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %v", err)
// 		}

// 		videos, err := ListVideos(db, ListVideosRequest{
// 			Limit:           *req.Limit,
// 			Offset:          *req.Offset,
// 			PlaylistID:      &playlist.ID,
// 			IncludeArchived: true,
// 			IncludeDrafts:   true,
// 		})
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to list videos: %v", err)
// 		}

// 		playlist.Videos = videos

// 		playlists = append(playlists, playlist)
// 	}

// 	return &playlists, nil
// }
