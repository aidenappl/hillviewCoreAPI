package query

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type UpdateCheckoutRequest struct {
	ID      *int                   `json:"id"`
	Changes *UpdateCheckoutChanges `json:"changes"`
}

type UpdateCheckoutChanges struct {
	// mutators
	CheckIn *bool `json:"check_in"`
}

func UpdateCheckout(db db.Queryable, req UpdateCheckoutRequest) (*structs.Checkout, error) {
	// validate fields
	if req.ID == nil {
		return nil, fmt.Errorf("no id provided")
	}

	// validate required fields
	if req.Changes == nil {
		return nil, fmt.Errorf("no changes provided")
	}

	if req.Changes.CheckIn == nil {
		return nil, fmt.Errorf("no changes provided")
	}

	if req.Changes.CheckIn != nil && *req.Changes.CheckIn {
		// check in the user
		query, args, err := sq.Update("asset_checkouts").
			Set("checkout_status", 2).
			Set("time_in", time.Unix(time.Now().Unix(), 0)).
			Set("checkout_notes", "checked in by admin").
			Where(sq.Eq{"id": req.ID}).
			ToSql()
		if err != nil {
			return nil, fmt.Errorf("error building query: %v", err)
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			return nil, fmt.Errorf("error running query: %v", err)
		}
	}

	// get the checkout
	return GetCheckout(db, GetCheckoutRequest{
		ID: req.ID,
	})

}
