package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListAssetsRequest struct {
	ID         *int
	Identifier *string
	UseOr      bool
	Limit      int
	Offset     int
	Search     *string
	Sort       *string
}

func ListAssets(db db.Queryable, req ListAssetsRequest) (*[]structs.Asset, error) {
	// check request
	if req.Limit == 0 {
		return nil, fmt.Errorf("no limit provided")
	}

	if req.Sort != nil {
		if *req.Sort != "DESC" && *req.Sort != "ASC" {
			return nil, fmt.Errorf("invalid sort direction provided")
		}
	} else {
		req.Sort = new(string)
		*req.Sort = "ASC"
	}

	// build query
	q := sq.Select(
		"assets.id",
		"assets.name",
		"assets.identifier",
		"assets.image_url",
		"assets.description",
		"assets.inserted_at",

		"asset_statuses.id",
		"asset_statuses.name",
		"asset_statuses.short_name",

		"asset_categories.id",
		"asset_categories.name",
		"asset_categories.short_name",

		"asset_metadata.serial_number",
		"asset_metadata.manufacturer",
		"asset_metadata.model",
		"asset_metadata.notes",
	).
		From("assets").
		Join("asset_statuses ON assets.status = asset_statuses.id").
		Join("asset_categories ON assets.category = asset_categories.id").
		LeftJoin("asset_metadata ON assets.id = asset_metadata.asset_id").
		OrderBy("assets.id " + *req.Sort).
		Limit(uint64(req.Limit)).
		Offset(uint64(req.Offset))

	if req.UseOr {
		q = q.Where(sq.Or{
			sq.Eq{"assets.id": req.ID},
			sq.Eq{"assets.identifier": req.Identifier},
			sq.Like{"assets.name": req.Search},
		})
	} else {
		if req.ID != nil {
			q = q.Where(sq.Eq{"assets.id": req.ID})
		}

		if req.Identifier != nil {
			q = q.Where(sq.Eq{"assets.identifier": req.Identifier})
		}

		if req.Search != nil {
			q = q.Where(sq.Like{"assets.name": req.Search})
		}
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}

	assets := []structs.Asset{}
	for rows.Next() {
		asset := structs.Asset{}
		assetStatus := structs.GeneralNSN{}
		assetCategory := structs.GeneralNSN{}
		assetMetadata := structs.AssetMetadata{}
		err := rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Identifier,
			&asset.ImageURL,
			&asset.Description,
			&asset.InsertedAt,

			&assetStatus.ID,
			&assetStatus.Name,
			&assetStatus.ShortName,

			&assetCategory.ID,
			&assetCategory.Name,
			&assetCategory.ShortName,

			&assetMetadata.SerialNumber,
			&assetMetadata.Manufacturer,
			&assetMetadata.Model,
			&assetMetadata.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		asset.Status = &assetStatus
		asset.Category = &assetCategory
		asset.Metadata = &assetMetadata

		assets = append(assets, asset)
	}

	return &assets, nil
}
