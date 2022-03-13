package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
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

type ListLinksRequest struct {
	Limit *uint64 `json:"limit"`
}

func ListLinks(db db.Queryable, req ListLinksRequest) ([]*structs.Link, error) {
	query, args, err := sq.Select(
		"links.id",
		"links.route",
		"links.destination",
		"links.active",
		"links.inserted_at",

		"users.id",
		"users.name",
		"users.email",
		"users.profile_image_url",
	).
		From("links").
		OrderBy("links.id").
		Join("users ON links.created_by = users.id").
		Limit(*req.Limit).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run sql query: %w", err)
	}

	defer rows.Close()

	links := []*structs.Link{}

	for rows.Next() {
		link := structs.Link{}
		creator := structs.UserTS{}
		err = rows.Scan(
			&link.ID,
			&link.Route,
			&link.Destination,
			&link.Active,
			&link.InsertedAt,

			&creator.ID,
			&creator.Name,
			&creator.Email,
			&creator.ProfileImageURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		link.Creator = creator

		links = append(links, &link)
	}

	return links, nil
}
