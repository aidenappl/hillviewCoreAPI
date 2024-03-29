package query

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type FindUserRequest struct {
	ID       *int    `json:"id"`
	Email    *string `json:"email"`
	Username *string `json:"username"`
}

func FindUser(db db.Queryable, req FindUserRequest) (*structs.User, error) {

	q := sq.Select(
		"users.id",
		"users.name",
		"users.username",
		"users.email",
		"users.profile_image_url",
		"users.inserted_at",

		"user_types.id",
		"user_types.name",
		"user_types.short_name",

		"user_authentication.google_id",
		"user_authentication.password",
	).
		From("users").
		Join("user_authentication ON user_authentication.user_id = users.id").
		Join("user_types ON user_types.id = users.authentication")

	if req.ID != nil {
		q = q.Where(sq.Eq{"users.id": req.ID})
	}

	if req.Email != nil {
		q = q.Where(sq.Eq{"users.email": req.Email})
	}

	if req.Username != nil {
		q = q.Where(sq.Eq{"users.username": req.Username})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	user := structs.User{}
	status := structs.GeneralNSN{}
	authentication := structs.UserAuthenticationStrategies{}
	err = rows.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.ProfileImageURL,
		&user.InsertedAt,

		&status.ID,
		&status.Name,
		&status.ShortName,

		&authentication.GoogleID,
		&authentication.Password,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning row: %w", err)
	}

	user.Authentication = status
	user.AuthenticationStrategies = &authentication

	return &user, nil
}

type ListAdminUsersRequest struct {
	Limit *uint64 `json:"limit"`
}

func ListAdminUsers(db db.Queryable, req ListAdminUsersRequest) ([]*structs.User, error) {

	if req.Limit == nil {
		return nil, fmt.Errorf("missing limit")
	}

	query, args, err := sq.Select(
		"users.id",
		"users.name",
		"users.username",
		"users.email",
		"users.profile_image_url",
		"users.inserted_at",

		"user_types.id",
		"user_types.name",
		"user_types.short_name",
	).
		From("users").
		Join("user_types ON user_types.id = users.authentication").
		Where(sq.NotEq{"user_types.id": 9}).
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

	users := []*structs.User{}

	for rows.Next() {

		user := structs.User{}
		status := structs.GeneralNSN{}
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.ProfileImageURL,
			&user.InsertedAt,

			&status.ID,
			&status.Name,
			&status.ShortName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sql rows: %w", err)
		}

		lastActive, err := GetUserLastActive(db, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get last active timestamp: %w", err)
		}

		user.LastActive = lastActive

		user.Authentication = status

		users = append(users, &user)

	}

	return users, nil

}

//
// Edit Admin User
//

type EditAdminAccountRequest struct {
	ID            int                       `json:"id"`
	Modifications AdminAccountModifications `json:"modifications"`
}

type AdminAccountModifications struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Authentication  *int    `json:"authentication"`
	ProfileImageURL *string `json:"profile_image_url"`
	Username        *string `json:"username"`
}

func EditAdminAccount(db db.Queryable, req EditAdminAccountRequest) (*structs.User, error) {
	q := sq.Update("users").
		Where(sq.Eq{"id": req.ID})

	if req.Modifications.Name != nil {
		q = q.Set("name", req.Modifications.Name)
	}

	if req.Modifications.Email != nil {
		q = q.Set("email", req.Modifications.Email)
	}

	if req.Modifications.Authentication != nil {
		q = q.Set("authentication", req.Modifications.Authentication)
	}

	if req.Modifications.ProfileImageURL != nil {
		q = q.Set("profile_image_url", req.Modifications.ProfileImageURL)
	}

	if req.Modifications.Username != nil {
		q = q.Set("username", req.Modifications.Username)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run sql query: %w", err)
	}

	return nil, nil

}

func InsertRequestLog(db db.Queryable, userID int, route string, method string) error {
	query, args, err := sq.Insert("request_logs").
		Columns(
			"user_id",
			"route",
			"method",
		).
		Values(
			userID,
			route,
			method,
		).ToSql()

	if err != nil {
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute sql query: %w", err)
	}

	return nil
}

func GetUserLastActive(db db.Queryable, userID int) (*time.Time, error) {
	query, args, err := sq.Select(
		"MAX(request_logs.inserted_at) as last_active",
	).
		From("request_logs").
		Where(sq.Eq{"request_logs.user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to successful execute query: %w", err)
	}

	defer rows.Close()

	// var timestamp time.Time

	var time *time.Time

	for rows.Next() {
		err = rows.Scan(
			&time,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sql rows: %w", err)
		}
	}

	return time, nil

}

func ReadUser(db db.Queryable, id *int, tag *string) (*structs.MobileUser, error) {

	if id == nil && tag == nil {
		return nil, fmt.Errorf("must specify either id or tag")
	}

	q := sq.Select(
		"users.id",
		"users.name",
		"users.email",
		"users.identifier",
		"users.profile_image_url",
		"users.inserted_at",
	).
		From("users")

	if id != nil {
		q = q.Where(sq.Eq{"users.id": *id})
	}

	if tag != nil {
		q = q.Where(sq.Eq{"users.identifier": *tag})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}

	if !rows.Next() {
		return nil, nil
	}

	defer rows.Close()

	var user structs.MobileUser

	err = rows.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Identifier,
		&user.ProfileImageURL,
		&user.InsertedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user row: %w", err)
	}

	return &user, nil

}
