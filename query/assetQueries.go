package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type EditAssetRequest struct {
	ID            *int                `json:"id"`
	Modifications *AssetModifications `json:"modifications"`
}

type AssetModifications struct {
	Name        *string `json:"name"`
	ImageURL    *string `json:"image_url"`
	Identifier  *string `json:"identifier"`
	Description *string `json:"description"`
	Status      *int    `json:"status"`
	Category    *int    `json:"category"`
	Notes       *string `json:"notes"`
}

type ListOpenCheckoutsRequest struct {
	Limit *uint64 `json:"limit"`
}

func ListOpenCheckouts(db db.Queryable, req ListOpenCheckoutsRequest) ([]*structs.Checkout, error) {
	query, args, err := sq.Select(
		"asset_checkouts.id",
		"asset_checkouts.asset_id",
		"asset_checkouts.offsite",
		"asset_checkouts.associated_user",
		"asset_checkouts.checkout_notes",

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

		"asset_checkouts.time_out",
		"asset_checkouts.time_in",
		"asset_checkouts.expected_in",
	).
		From("asset_checkouts").
		LeftJoin("checkout_statuses ON asset_checkouts.checkout_status = checkout_statuses.id").
		LeftJoin("assets ON asset_checkouts.asset_id = assets.id").
		LeftJoin("users ON asset_checkouts.associated_user = users.id").
		OrderBy("asset_checkouts.id DESC").
		Where("asset_checkouts.checkout_status = 1").
		Limit(*req.Limit).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
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

			&checkout.TimeOut,
			&checkout.TimeIn,
			&checkout.ExpectedIn,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if checkout.ExpectedIn.IsZero() {
			checkout.ExpectedIn = nil
		}

		checkout.Asset = &asset

		checkout.User = &user

		checkout.CheckoutStatus = &checkout_status

		checkouts = append(checkouts, &checkout)
	}

	return checkouts, nil
}

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

		"asset_checkouts.time_out",
		"asset_checkouts.time_in",
		"asset_checkouts.expected_in",
	).
		From("asset_checkouts").
		LeftJoin("checkout_statuses ON asset_checkouts.checkout_status = checkout_statuses.id").
		LeftJoin("assets ON asset_checkouts.asset_id = assets.id").
		LeftJoin("users ON asset_checkouts.associated_user = users.id").
		OrderBy("asset_checkouts.id DESC").
		Limit(*req.Limit).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
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

			&checkout.TimeOut,
			&checkout.TimeIn,
			&checkout.ExpectedIn,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if checkout.ExpectedIn.IsZero() {
			checkout.ExpectedIn = nil
		}

		checkout.Asset = &asset

		checkout.User = &user

		checkout.CheckoutStatus = &checkout_status

		checkouts = append(checkouts, &checkout)
	}

	return checkouts, nil
}
