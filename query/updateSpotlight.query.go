package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type UpdateSpotlightRequest struct {
	Rank    *int `json:"rank"`
	VideoID *int `json:"video_id"`
}

func UpdateSpotlight(db db.Queryable, req UpdateSpotlightRequest) (*structs.Spotlight, error) {
	//  check required fields
	if req.Rank == nil {
		return nil, fmt.Errorf("required field rank is nil")
	}

	if req.VideoID == nil {
		return nil, fmt.Errorf("required field video_id is nil")
	}

	// build query
	q := sq.Update("spotlight").
		Set("video_id", *req.VideoID).
		Where(sq.Eq{"rank": *req.Rank})

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %v", err)
	}

	// execute query
	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// get updated spotlight
	spotlight, err := GetSpotlight(db, GetSpotlightRequest{Rank: req.Rank})
	if err != nil {
		return nil, fmt.Errorf("error getting updated spotlight: %v", err)
	}

	return spotlight, nil
}
