package query

import (
	"fmt"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type GetCheckoutRequest struct {
	ID *int
}

func GetCheckout(db db.Queryable, req GetCheckoutRequest) (*structs.Checkout, error) {
	// validate fields
	if req.ID == nil {
		return nil, fmt.Errorf("no id provided")
	}

	// use the list req
	checkouts, err := ListCheckouts(db, ListCheckoutsRequest{
		ID:     req.ID,
		Limit:  &[]int{1}[0],
		Offset: &[]int{0}[0],
	})
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	if len(checkouts) == 0 {
		return nil, nil
	}

	return checkouts[0], nil
}
