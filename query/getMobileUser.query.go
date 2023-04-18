package query

import (
	"fmt"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type GetMobileUserRequest struct {
	ID         *int
	Identifier *string
}

func GetMobileUser(db db.Queryable, req GetMobileUserRequest) (*structs.MobileUser, error) {
	// validate required fields
	if req.ID == nil && req.Identifier == nil {
		return nil, fmt.Errorf("no id or identifier provided")
	}

	// use the list req
	users, err := ListMobileUsers(db, ListMobileUsersRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
		Limit:      &[]int{1}[0],
		Offset:     &[]int{0}[0],
	})
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users[0], nil
}
