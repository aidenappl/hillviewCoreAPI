package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListAssetsRequest struct {
	Limit *uint64 `json:"limit"`
}

func ListAssets(db db.Queryable, req ListAssetsRequest) ([]*structs.Asset, error) {

	q := sq.Select(
		"assets.id",
		"assets.name",
		"assets.identifier",
		"assets.image_url",
		"assets.description",

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
		Limit(*req.Limit)

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

	assets := []*structs.Asset{}

	for rows.Next() {
		asset := structs.Asset{}
		assetStatus := structs.GeneralNSN{}
		assetCategory := structs.GeneralNSN{}
		assetMetadata := structs.AssetMetadata{}

		err = rows.Scan(
			&asset.ID,
			&asset.Name,
			&asset.Identifier,
			&asset.ImageURL,
			&asset.Description,

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

		checkout, err := ReadActiveCheckouts(db, int(asset.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to read active checkouts: %w", err)
		}

		if checkout != nil {
			user, err := ReadUser(db, checkout.AssociatedUser, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to read active checkouts: %w", err)
			}
			checkout.User = user
		}

		asset.ActiveTab = checkout

		assets = append(assets, &asset)
	}

	return assets, nil
}
