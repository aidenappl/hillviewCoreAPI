package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type UpdateAssetRequest struct {
	// Ident Fields from Query
	ID         *int    `json:"id"`
	Identifier *string `json:"identifier"`

	// Edit Fields
	Changes *UpdateAssetChanges `json:"changes"`
}

type UpdateAssetChanges struct {
	// Asset Fields
	Name        *string                `json:"name"`
	ImageURL    *string                `json:"image_url"`
	Identifier  *string                `json:"identifier"`
	Description *string                `json:"description"`
	Category    *int                   `json:"category"`
	Status      *int                   `json:"status"`
	Metadata    *structs.AssetMetadata `json:"metadata"`
}

func UpdateAsset(db db.Queryable, req UpdateAssetRequest) (*structs.Asset, error) {
	// check required fields
	if req.ID == nil && req.Identifier == nil {
		return nil, fmt.Errorf("id or identifier is required")
	}

	// check changes
	if req.Changes == nil {
		return nil, fmt.Errorf("changes is required")
	}

	if req.Changes.Metadata != nil {
		if req.ID == nil {
			return nil, fmt.Errorf("id is required to update metadata")
		}

		// run metadata update
		_, err := UpdateAssetMetadata(db, UpdateAssetMetadataRequest{
			ID:      req.ID,
			Changes: req.Changes.Metadata,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update metadata: %v", err)
		}
	}
	skipUpdate := false
	if req.Changes.Metadata != nil && req.Changes.Name == nil && req.Changes.ImageURL == nil && req.Changes.Identifier == nil && req.Changes.Description == nil && req.Changes.Category == nil && req.Changes.Status == nil {
		skipUpdate = true
	}

	// check if metadata is the only change
	if !skipUpdate {
		// build query
		q := sq.Update("assets")

		// set changes
		if req.Changes.Name != nil {
			q = q.Set("name", *req.Changes.Name)
		}

		if req.Changes.ImageURL != nil {
			q = q.Set("image_url", *req.Changes.ImageURL)
		}

		if req.Changes.Identifier != nil {
			q = q.Set("identifier", *req.Changes.Identifier)
		}

		if req.Changes.Description != nil {
			q = q.Set("description", *req.Changes.Description)
		}

		if req.Changes.Category != nil {
			q = q.Set("category", *req.Changes.Category)
		}

		if req.Changes.Status != nil {
			q = q.Set("status", *req.Changes.Status)
		}

		// set where
		if req.ID != nil {
			q = q.Where(sq.Eq{"id": *req.ID})
		}

		if req.Identifier != nil {
			q = q.Where(sq.Eq{"identifier": *req.Identifier})
		}

		// build query
		query, args, err := q.ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build query: %v", err)
		}

		// execute query
		_, err = db.Exec(query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %v", err)
		}
	}

	return GetAsset(db, GetAssetRequest{
		ID: req.ID, Identifier: req.Identifier})

}
