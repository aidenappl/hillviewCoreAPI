package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type UpdateAssetMetadataRequest struct {
	ID      *int                   `json:"id"`
	Changes *structs.AssetMetadata `json:"changes"`
}

func UpdateAssetMetadata(db db.Queryable, req UpdateAssetMetadataRequest) (*structs.AssetMetadata, error) {
	// check required fields
	if req.ID == nil {
		return nil, fmt.Errorf("id is required")
	}

	// check changes
	if req.Changes == nil {
		return nil, fmt.Errorf("changes is required")
	}

	// build query
	q := sq.Update("asset_metadata")

	// set changes
	if req.Changes.SerialNumber != nil {
		q = q.Set("serial_number", *req.Changes.SerialNumber)
	}

	if req.Changes.Manufacturer != nil {
		q = q.Set("manufacturer", *req.Changes.Manufacturer)
	}

	if req.Changes.Model != nil {
		q = q.Set("model", *req.Changes.Model)
	}

	if req.Changes.Notes != nil {
		q = q.Set("notes", *req.Changes.Notes)
	}

	// set where
	q = q.Where(sq.Eq{"asset_id": *req.ID})

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

	// return asset
	return GetAssetMetadata(db, GetAssetMetadataRequest{
		ID: req.ID,
	})
}
