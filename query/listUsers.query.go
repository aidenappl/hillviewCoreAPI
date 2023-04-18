package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListUsersRequest struct {
	Limit  *int
	Offset *int

	// reqs
	Search               *string
	Sort                 *string
	ID                   *int
	IncludeSensitiveData bool
}

func ListUsers(db db.Queryable, req ListUsersRequest) ([]*structs.User, error) {
	// check required params
	if req.Limit == nil {
		return nil, fmt.Errorf("limit is required")
	}

	if req.Offset == nil {
		return nil, fmt.Errorf("offset is required")
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
		"users.username",
		"users.name",
		"users.email",
		"users.profile_image_url",
		"users.inserted_at",
		`(SELECT MAX(request_logs.inserted_at)
        FROM request_logs
        WHERE request_logs.user_id = users.id) as last_active`,

		"user_types.id",
		"user_types.name",
		"user_types.short_name",

		"user_authentication.google_id",
		"user_authentication.password",
	).
		From("users").
		LeftJoin("user_types ON users.authentication = user_types.id").
		LeftJoin("user_authentication ON users.id = user_authentication.user_id").
		OrderBy("users.id " + *req.Sort).
		Limit(uint64(*req.Limit)).
		Offset(uint64(*req.Offset))

	// add search
	if req.Search != nil {
		q = q.Where(sq.Or{
			sq.Like{"users.username": *req.Search},
			sq.Like{"users.name": *req.Search},
			sq.Like{"users.email": *req.Search},
		})
	}

	// add id
	if req.ID != nil {
		q = q.Where(sq.Eq{"users.id": *req.ID})
	}

	// run the query
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	// execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	// scan the results
	var users []*structs.User
	for rows.Next() {
		var user structs.User
		var userAuth structs.UserAuthenticationStrategies
		var userType structs.GeneralNSN
		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Name,
			&user.Email,
			&user.ProfileImageURL,
			&user.InsertedAt,
			&user.LastActive,

			&userType.ID,
			&userType.Name,
			&userType.ShortName,

			&userAuth.GoogleID,
			&userAuth.Password,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if req.IncludeSensitiveData {
			user.AuthenticationStrategies = &userAuth
		}

		user.Authentication = userType

		users = append(users, &user)
	}

	return users, nil
}
