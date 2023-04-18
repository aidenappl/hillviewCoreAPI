package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type UpdateMobileUserRequest struct {
	ID      *int                     `json:"id"`
	Changes *UpdateMobileUserChanges `json:"changes"`
}

type UpdateMobileUserChanges struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Identifier      *string `json:"identifier"`
	Status          *int    `json:"status"`
	ProfileImageURL *string `json:"profile_image_url"`
}

func UpdateMobileUser(db db.Queryable, req UpdateMobileUserRequest) (*structs.MobileUser, error) {
	// validate the fields
	if req.ID == nil {
		return nil, fmt.Errorf("id is required")
	}

	// validate changes
	if req.Changes == nil {
		return nil, fmt.Errorf("changes is required")
	}

	// check that modifications exist
	if req.Changes.Name == nil &&
		req.Changes.Email == nil &&
		req.Changes.Identifier == nil &&
		req.Changes.Status == nil &&
		req.Changes.ProfileImageURL == nil {
		return nil, fmt.Errorf("no modifications were made")
	}

	// build the query
	q := sq.Update("users").
		Where(sq.Eq{"id": req.ID})

	if req.Changes.Name != nil {
		q = q.Set("name", req.Changes.Name)
	}

	if req.Changes.Email != nil {
		q = q.Set("email", req.Changes.Email)
	}

	if req.Changes.Identifier != nil {
		q = q.Set("identifier", req.Changes.Identifier)
	}

	if req.Changes.Status != nil {
		q = q.Set("status", req.Changes.Status)
	}

	if req.Changes.ProfileImageURL != nil {
		q = q.Set("profile_image_url", req.Changes.ProfileImageURL)
	}

	// execute the query
	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// get the updated user
	return GetMobileUser(db, GetMobileUserRequest{
		ID: req.ID,
	})
}
