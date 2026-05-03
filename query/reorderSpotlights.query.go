package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ReorderSpotlightItem struct {
	Position int `json:"position"`
	VideoID  int `json:"video_id"`
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
			Where(sq.Eq{"position": item.Position})

		query, args, err := q.ToSql()
		if err != nil {
			return nil, fmt.Errorf("error building query for position %d: %v", item.Position, err)
		}

		_, err = database.Exec(query, args...)
		if err != nil {
			return nil, fmt.Errorf("error executing update for position %d: %v", item.Position, err)
		}
	}

	limit := 50
	offset := 0
	return ListSpotlights(database, ListSpotlightsRequest{Limit: &limit, Offset: &offset})
}
