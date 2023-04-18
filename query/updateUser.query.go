package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type UpdateUserRequest struct {
	ID      *int               `json:"id"`
	Changes *UpdateUserChanges `json:"changes"`
}

type UpdateUserChanges struct {
	Username        *string `json:"username"`
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	ProfileImageURL *string `json:"profile_image_url"`
	Authentication  *int    `json:"authentication"`
}

func UpdateUser(db db.Queryable, req UpdateUserRequest) (*structs.User, error) {
	// validate the fields
	if req.ID == nil {
		return nil, fmt.Errorf("id is required")
	}

	// validate changes
	if req.Changes == nil {
		return nil, fmt.Errorf("changes is required")
	}

	// check that modifications exist
	if req.Changes.Username == nil &&
		req.Changes.Name == nil &&
		req.Changes.Email == nil &&
		req.Changes.ProfileImageURL == nil &&
		req.Changes.Authentication == nil {
		return nil, fmt.Errorf("no modifications were made")
	}

	// build the query
	q := sq.Update("users").
		Where(sq.Eq{"id": req.ID})

	if req.Changes.Username != nil {
		q = q.Set("username", req.Changes.Username)
	}

	if req.Changes.Name != nil {
		q = q.Set("name", req.Changes.Name)
	}

	if req.Changes.Email != nil {
		q = q.Set("email", req.Changes.Email)
	}

	if req.Changes.ProfileImageURL != nil {
		q = q.Set("profile_image_url", req.Changes.ProfileImageURL)
	}

	if req.Changes.Authentication != nil {
		q = q.Set("authentication", req.Changes.Authentication)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to run sql query: %w", err)
	}

	return GetUser(db, GetUserRequest{
		ID: req.ID,
	})
}
