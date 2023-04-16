package query

import (
	"fmt"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type GetPlaylistRequest struct {
	ID         *int
	Identifier *string
}

func GetPlaylist(db db.Queryable, req GetPlaylistRequest) (*structs.Playlist, error) {
	// validate the request
	if req.ID == nil && req.Identifier == nil {
		return nil, fmt.Errorf("no id or identifier provided")
	}

	// run list query to get the asset
	playlists, err := ListPlaylists(db, ListPlaylistsRequest{
		ID:     req.ID,
		Limit:  &[]int{1}[0],
		Offset: &[]int{0}[0],
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// check if the asset exists
	if len(playlists) == 0 {
		return nil, nil
	}

	return playlists[0], nil

}
