package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
)

type CreateLinkRequest struct {
	Route    *string `json:"route"`
	Endpoint *string `json:"endpoint"`
	Creator  *int    `json:"user"`
}

func CreateLink(db db.Queryable, req CreateLinkRequest) error {
	query, args, err := sq.Insert("links").
		Columns("route", "destination", "created_by").
		Values(
			req.Route,
			req.Endpoint,
			req.Creator,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to run sql query: %w", err)
	}

	return nil
}
