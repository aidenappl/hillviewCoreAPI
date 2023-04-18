package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type MobileAccountModifications struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	ProfileImageURL *string `json:"profile_image_url"`
	Identifier      *string `json:"identifier"`
	Status          *int    `json:"status"`
}

type EditMobileAccountRequest struct {
	ID            *int                        `json:"id"`
	Modifications *MobileAccountModifications `json:"modifications"`
}

func EditMobileAccount(db db.Queryable, req EditMobileAccountRequest) (interface{}, error) {
	q := sq.Update("users").
		Where(sq.Eq{"users.id": req.ID})

	if req.Modifications.Name != nil {
		q = q.Set("users.name", req.Modifications.Name)
	}
	if req.Modifications.Email != nil {
		q = q.Set("users.email", req.Modifications.Email)
	}
	if req.Modifications.ProfileImageURL != nil {
		q = q.Set("users.profile_image_url", req.Modifications.ProfileImageURL)
	}
	if req.Modifications.Identifier != nil {
		q = q.Set("users.identifier", req.Modifications.Identifier)
	}
	if req.Modifications.Status != nil {
		q = q.Set("users.status", req.Modifications.Status)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run sql query: %w", err)
	}

	return GetMobileUsers(db, GetMobileUsersRequest{ID: req.ID})

}

type GetMobileUsersRequest struct {
	ID     *int
	Limit  *uint64
	Offset *uint64
	Search *string
	Sort   *string
}

func GetMobileUsers(db db.Queryable, req GetMobileUsersRequest) (interface{}, error) {
	q := sq.Select(
		"users.id",
		"users.name",
		"users.identifier",
		"users.email",
		"users.profile_image_url",
		"users.inserted_at",

		"user_statuses.id",
		"user_statuses.name",
		"user_statuses.short_name",
	).
		From("users").
		Join("user_statuses ON users.status = user_statuses.id")

	if req.ID != nil && req.Search != nil {
		return nil, fmt.Errorf("cannot search and pass ID")
	}

	if req.ID != nil {
		q = q.Where(sq.Eq{"users.id": req.ID})
	}

	if req.Limit != nil && req.Offset != nil && req.Sort != nil {
		q = q.Limit(*req.Limit).Offset(*req.Offset).OrderBy("users.id " + *req.Sort)
	}

	if req.Search != nil {
		q = q.Where(sq.Or{
			sq.Like{"users.name": "%" + *req.Search + "%"},
			sq.Like{"users.identifier": "%" + *req.Search + "%"},
			sq.Like{"users.email": "%" + *req.Search + "%"},
		})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	if req.ID != nil {
		user := structs.MobileUser{}
		userStatus := structs.GeneralNSN{}
		err := db.QueryRow(query, args...).Scan(
			&user.ID,
			&user.Name,
			&user.Identifier,
			&user.Email,
			&user.ProfileImageURL,
			&user.InsertedAt,

			&userStatus.ID,
			&userStatus.Name,
			&userStatus.ShortName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sql rows: %w", err)
		}
		user.Status = &userStatus
		return &user, nil
	} else {
		rows, err := db.Query(query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to run sql query: %w", err)
		}

		defer rows.Close()

		users := []*structs.MobileUser{}

		for rows.Next() {

			user := structs.MobileUser{}
			userStatus := structs.GeneralNSN{}
			err := rows.Scan(
				&user.ID,
				&user.Name,
				&user.Identifier,
				&user.Email,
				&user.ProfileImageURL,
				&user.InsertedAt,

				&userStatus.ID,
				&userStatus.Name,
				&userStatus.ShortName,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan sql rows: %w", err)
			}

			user.Status = &userStatus

			users = append(users, &user)
		}
		return users, nil
	}
	return nil, nil
}
