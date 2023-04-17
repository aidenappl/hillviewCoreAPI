package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type CreateAssetRequest struct {
	Name        *string                `json:"name"`
	ImageURL    *string                `json:"image_url"`
	Identifier  *string                `json:"identifier"`
	Description *string                `json:"description"`
	Category    *int                   `json:"category"`
	Metadata    *structs.AssetMetadata `json:"metadata"`
}

func CreateAsset(db db.Queryable, req CreateAssetRequest) (*structs.Asset, error) {
	// validate fields
	if req.Name == nil || *req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if req.Identifier == nil || *req.Identifier == "" {
		return nil, fmt.Errorf("identifier is required")
	}

	if req.Category == nil {
		return nil, fmt.Errorf("category is required")
	}

	if req.ImageURL == nil || *req.ImageURL == "" {
		return nil, fmt.Errorf("image_url is required")
	}

	if req.Description == nil || *req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	// validate metadata
	if req.Metadata == nil {
		return nil, fmt.Errorf("metadata is required")
	}

	if req.Metadata.Manufacturer == nil || *req.Metadata.Manufacturer == "" {
		return nil, fmt.Errorf("metadata.manufacturer is required")
	}

	if req.Metadata.Model == nil || *req.Metadata.Model == "" {
		return nil, fmt.Errorf("metadata.model is required")
	}

	// insert asset
	query, args, err := sq.Insert("assets").
		Columns(
			"name",
			"image_url",
			"identifier",
			"description",
			"category",
		).
		Values(
			*req.Name,
			*req.ImageURL,
			*req.Identifier,
			*req.Description,
			*req.Category,
		).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create insert query: %v", err)
	}

	// execute query
	rows, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert asset: %v", err)
	}

	// get asset id
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %v", err)
	}

	// insert metadata
	query, args, err = sq.Insert("asset_metadata").
		Columns(
			"asset_id",
			"manufacturer",
			"model",
			"serial_number",
		).
		Values(
			id,
			*req.Metadata.Manufacturer,
			*req.Metadata.Model,
			*req.Metadata.SerialNumber,
		).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create insert query: %v", err)
	}

	// execute query
	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert asset metadata: %v", err)
	}

	// convert id to int
	idInt := int(id)

	// get asset
	return GetAsset(db, GetAssetRequest{
		ID: &idInt,
	})
}
