package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ReorderSpotlightItem struct {
	Rank    int `json:"rank"`
	VideoID int `json:"video_id"`
}

type ReorderSpotlightsRequest struct {
	Items []ReorderSpotlightItem
}

func ReorderSpotlights(database db.Queryable, req ReorderSpotlightsRequest) ([]*structs.Spotlight, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("no items provided")
	}

	for _, item := range req.Items {
		q := sq.Update("spotlight").
			Set("video_id", item.VideoID).
			Where(sq.Eq{"rank": item.Rank})

		query, args, err := q.ToSql()
		if err != nil {
			return nil, fmt.Errorf("error building query for rank %d: %v", item.Rank, err)
		}

		_, err = database.Exec(query, args...)
		if err != nil {
			return nil, fmt.Errorf("error executing update for rank %d: %v", item.Rank, err)
		}
	}

	limit := 50
	offset := 0
	return ListSpotlights(database, ListSpotlightsRequest{Limit: &limit, Offset: &offset})
}
