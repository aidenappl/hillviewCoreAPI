package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type ListCheckoutsRequest struct {
	Limit *uint64 `json:"limit"`
}

func ListCheckouts(db db.Queryable, req ListCheckoutsRequest) ([]*structs.Checkout, error) {
	query, args, err := sq.Select(
		"asset_checkouts.id",
		"asset_checkouts.asset_id",
		"asset_checkouts.offsite",
		"asset_checkouts.associated_user",
		"asset_checkouts.checkout_notes",

		"asset_checkouts.time_out",
		"asset_checkouts.time_in",
		"asset_checkouts.expected_in",

		"checkout_statuses.id",
		"checkout_statuses.name",
		"checkout_statuses.short_name",

		"users.id",
		"users.name",
		"users.email",
		"users.identifier",
		"users.profile_image_url",
		"users.inserted_at",

		"assets.id",
		"assets.name",
		"assets.identifier",
		"assets.image_url",
		"assets.description",
	).
		From("asset_checkouts").
		LeftJoin("checkout_statuses ON asset_checkouts.checkout_status = checkout_statuses.id").
		LeftJoin("assets ON asset_checkouts.asset_id = assets.id").
		LeftJoin("users ON asset_checkouts.associated_user = users.id").
		Limit(*req.Limit).
		ToSql()

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

	checkouts := []*structs.Checkout{}

	for rows.Next() {
		var checkout structs.Checkout
		var checkout_status structs.GeneralNSN
		var user structs.MobileUser
		var asset structs.Asset

		err = rows.Scan(
			&checkout.ID,
			&checkout.AssetID,
			&checkout.Offsite,
			&checkout.AssociatedUser,
			&checkout.CheckoutNotes,

			&checkout.TimeOut,
			&checkout.TimeIn,
			&checkout.ExpectedIn,

			&checkout_status.ID,
			&checkout_status.Name,
			&checkout_status.ShortName,

			&user.ID,
			&user.Name,
			&user.Email,
			&user.Identifier,
			&user.ProfileImageURL,
			&user.InsertedAt,

			&asset.ID,
			&asset.Name,
			&asset.Identifier,
			&asset.ImageURL,
			&asset.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		checkout.Asset = &asset

		checkout.User = &user

		checkout.CheckoutStatus = &checkout_status

		checkouts = append(checkouts, &checkout)
	}

	return checkouts, nil
}

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

func ReadAsset(db db.Queryable, id *int, tag *string) (*structs.Asset, error) {

	if id == nil && tag == nil {
		return nil, fmt.Errorf("must specify either id or tag")
	}

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
		LeftJoin("asset_metadata ON assets.id = asset_metadata.asset_id")

	if id != nil {
		q = q.Where(sq.Eq{"assets.id": id})
	}

	if tag != nil {
		q = q.Where(sq.Eq{"assets.identifier": tag})
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

	return &asset, nil
}
