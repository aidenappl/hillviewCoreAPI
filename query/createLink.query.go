package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type CreateLinkRequest struct {
	Route       *string `json:"route"`
	Destination *string `json:"destination"`
}

func CreateLink(db db.Queryable, req CreateLinkRequest) (*structs.Link, error) {
	// validate the request
	if req.Route == nil || len(*req.Route) == 0 || req.Destination == nil || len(*req.Destination) == 0 {
		return nil, fmt.Errorf("missing required fields: route, destination")
	}

	// create the link
	query, args, err := sq.Insert("links").
		Columns("route", "destination").
		Values(*req.Route, *req.Destination).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql request: %v", err)
	}

	// execute the query
	rows, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql request: %v", err)
	}

	// get row id
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last inserted id: %v", err)
	}

	// convert id to int
	linkID := int(id)

	// return
	return GetLink(db, GetLinkRequest{
		ID: &linkID,
	})
}
