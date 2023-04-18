package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type CreateMobileUserRequest struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Identifier      *string `json:"identifier"`
	ProfileImageURL *string `json:"profile_image_url"`
}

func CreateMobileUser(db db.Queryable, req CreateMobileUserRequest) (*structs.MobileUser, error) {
	// validate the request
	if req.Name == nil {
		return nil, fmt.Errorf("name is required")
	}

	if req.Email == nil {
		return nil, fmt.Errorf("email is required")
	}

	if req.Identifier == nil {
		return nil, fmt.Errorf("identifier is required")
	}

	// build the query
	cols := []string{"name", "email", "identifier"}
	vals := []interface{}{req.Name, req.Email, req.Identifier}

	if req.ProfileImageURL != nil {
		cols = append(cols, "profile_image_url")
		vals = append(vals, req.ProfileImageURL)
	}

	query, args, err := sq.Insert("users").
		Columns(cols...).
		Values(vals...).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// run the query
	rows, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}

	// get user id
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// convert id to int
	idInt := int(id)

	// build the response
	return GetMobileUser(db, GetMobileUserRequest{
		ID: &idInt,
	})

}
