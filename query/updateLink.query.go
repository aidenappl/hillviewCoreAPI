package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type UpdateLinkRequest struct {
	ID      *int
	Changes *UpdateLinkChanges
}

type UpdateLinkChanges struct {
	Route       *string
	Destination *string
	Active      *bool
}

func UpdateLink(db db.Queryable, req UpdateLinkRequest) (*structs.Link, error) {
	// check required fields
	if req.ID == nil {
		return nil, fmt.Errorf("required field id is nil")
	}

	if req.Changes == nil {
		return nil, fmt.Errorf("required field changes is nil")
	}

	// build query
	q := sq.Update("links").
		Where(sq.Eq{"id": *req.ID})

	// add changes
	if req.Changes.Route != nil {
		q = q.Set("route", *req.Changes.Route)
	}

	if req.Changes.Destination != nil {
		q = q.Set("destination", *req.Changes.Destination)
	}

	if req.Changes.Active != nil {
		q = q.Set("active", *req.Changes.Active)
	}

	// form query
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error forming query: %v", err)
	}

	// execute query
	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating link: %v", err)
	}

	// get updated playlist
	return GetLink(db, GetLinkRequest{ID: req.ID})
}
