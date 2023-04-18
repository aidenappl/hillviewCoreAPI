package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListMobileUsersRequest struct {
	Limit  *int
	Offset *int
	Sort   *string

	// reqs
	Search     *string
	ID         *int
	Identifier *string
}

func ListMobileUsers(db db.Queryable, req ListMobileUsersRequest) ([]*structs.MobileUser, error) {
	// validate fields
	if req.Limit == nil {
		return nil, fmt.Errorf("no limit provided")
	}

	if req.Offset == nil {
		return nil, fmt.Errorf("no offset provided")
	}

	// check sort formatting
	if req.Sort != nil {
		if *req.Sort != "DESC" && *req.Sort != "ASC" {
			return nil, fmt.Errorf("sort must be DESC or ASC")
		}
	} else {
		// default to DESC
		req.Sort = new(string)
		*req.Sort = "DESC"
	}

	// build the query
	q := sq.Select(
		"users.id",
		"users.name",
		"users.email",
		"users.identifier",
		"users.profile_image_url",
		"users.inserted_at",

		"user_statuses.id",
		"user_statuses.name",
		"user_statuses.short_name",
	).From("users").
		Join("user_statuses ON user_statuses.id = users.user_status_id").
		Limit(uint64(*req.Limit)).
		Offset(uint64(*req.Offset)).
		OrderBy("users.inserted_at " + *req.Sort)

	// add search
	if req.Search != nil {
		q = q.Where(sq.Or{
			sq.Like{"users.name": "%" + *req.Search + "%"},
			sq.Like{"users.email": "%" + *req.Search + "%"},
			sq.Like{"users.identifier": "%" + *req.Search + "%"},
		})
	}

	// add id
	if req.ID != nil {
		q = q.Where(sq.Eq{"users.id": *req.ID})
	}

	// add identifier
	if req.Identifier != nil {
		q = q.Where(sq.Eq{"users.identifier": *req.Identifier})
	}

	// run the query
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run sql query: %w", err)
	}

	// parse the results
	users := []*structs.MobileUser{}
	for rows.Next() {
		user := structs.MobileUser{}
		userStatus := structs.GeneralNSN{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Identifier,
			&user.ProfileImageURL,
			&user.InsertedAt,

			&userStatus.ID,
			&userStatus.Name,
			&userStatus.ShortName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		user.Status = userStatus

		users = append(users, &user)
	}

	return users, nil

}
