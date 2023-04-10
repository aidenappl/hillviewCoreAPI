package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

func ReadActiveCheckouts(db db.Queryable, id int) (*structs.Checkout, error) {

	query, args, err := sq.Select(
		"asset_checkouts.id",
		"asset_checkouts.asset_id",
		"asset_checkouts.associated_user",
		"asset_checkouts.checkout_notes",
		"asset_checkouts.time_out",
		"asset_checkouts.time_in",
		"asset_checkouts.expected_in",

		"checkout_statuses.id",
		"checkout_statuses.name",
		"checkout_statuses.short_name",
	).
		From("asset_checkouts").
		Where(sq.Eq{"asset_checkouts.asset_id": id}).
		Where(sq.Eq{"asset_checkouts.checkout_status": 1}).
		Join("checkout_statuses ON asset_checkouts.checkout_status = checkout_statuses.id").
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

	var checkout structs.Checkout
	var checkout_status structs.GeneralNSN

	err = rows.Scan(
		&checkout.ID,
		&checkout.AssetID,
		&checkout.AssociatedUser,
		&checkout.CheckoutNotes,
		&checkout.TimeOut,
		&checkout.TimeIn,
		&checkout.ExpectedIn,

		&checkout_status.ID,
		&checkout_status.Name,
		&checkout_status.ShortName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	checkout.CheckoutStatus = &checkout_status

	return &checkout, nil
}
