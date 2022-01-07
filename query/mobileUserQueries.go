package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListMobileUsersRequest struct {
	Limit *uint64 `json:"limit"`
}

func ListMobileUsers(db db.Queryable, req ListMobileUsersRequest) ([]*structs.MobileUser, error) {

	if req.Limit == nil {
		return nil, fmt.Errorf("missing limit")
	}

	query, args, err := sq.Select(
		"users.id",
		"users.name",
		"users.identifier",
		"users.email",
		"users.profile_image_url",
		"users.inserted_at",
	).
		From("users").
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

	users := []*structs.MobileUser{}

	for rows.Next() {

		user := structs.MobileUser{}
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Identifier,
			&user.Email,
			&user.ProfileImageURL,
			&user.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sql rows: %w", err)
		}

		users = append(users, &user)

	}

	return users, nil

}
