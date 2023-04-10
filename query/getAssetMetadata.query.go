package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type GetAssetMetadataRequest struct {
	// Ident Fields from Query
	ID *int `json:"id"`
}

func GetAssetMetadata(db db.Queryable, req GetAssetMetadataRequest) (*structs.AssetMetadata, error) {
	// check required fields
	if req.ID == nil {
		return nil, fmt.Errorf("id is required")
	}

	// build query
	query, args, err := sq.Select(
		"serial_number",
		"manufacturer",
		"model",
		"notes",
	).From("asset_metadata").Where(sq.Eq{"asset_id": *req.ID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	// execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	// scan rows
	var assetMetadata structs.AssetMetadata
	if rows.Next() {
		err = rows.Scan(
			&assetMetadata.SerialNumber,
			&assetMetadata.Manufacturer,
			&assetMetadata.Model,
			&assetMetadata.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
	}

	// return asset
	return &assetMetadata, nil
}
