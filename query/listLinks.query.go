package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListLinksRequest struct {
	Limit  *int
	Offset *int
	Search *string
	Sort   *string
	SortBy *string
	Active *bool // nil = all, true = active only, false = archived only

	// params
	ID *int
}

func ListLinks(db db.Queryable, req ListLinksRequest) ([]*structs.Link, error) {
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

	sortDir := "DESC"
	if req.Sort != nil && *req.Sort == "asc" {
		sortDir = "ASC"
	}

	// determine order column
	orderCol := "links.inserted_at"
	if req.SortBy != nil {
		switch *req.SortBy {
		case "clicks":
			orderCol = "clicks"
		case "date":
			orderCol = "links.inserted_at"
		}
	}

	// build query
	q := sq.Select(
		"links.id",
		"links.route",
		"links.destination",
		"links.active",
		"links.inserted_at",

		"users.id",
		"users.name",
		"users.email",
		"users.profile_image_url",

		`(
			SELECT COUNT(*) FROM link_clicks WHERE link_clicks.link_id = links.id
		) as clicks`,
	).
		From("links").
		OrderBy(fmt.Sprintf("%s %s", orderCol, sortDir)).
		Join("users ON links.created_by = users.id").
		Limit(uint64(*req.Limit)).
		Offset(uint64(*req.Offset))

	// active filter: nil = all, true = active only, false = archived only
	if req.Active != nil {
		q = q.Where(sq.Eq{"links.active": *req.Active})
	}

	// add params
	if req.ID != nil {
		q = q.Where(sq.Eq{"links.id": *req.ID})
	}

	// add search
	if req.Search != nil {
		q = q.Where(sq.Like{"links.route": "%" + *req.Search + "%"})
	}

	// run query
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run sql query: %w", err)
	}

	links := []*structs.Link{}
	for rows.Next() {
		var link structs.Link
		var user structs.UserTS

		err = rows.Scan(
			&link.ID,
			&link.Route,
			&link.Destination,
			&link.Active,
			&link.InsertedAt,

			&user.ID,
			&user.Name,
			&user.Email,
			&user.ProfileImageURL,

			&link.Clicks,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		link.Creator = user
		links = append(links, &link)
	}

	return links, nil
}
