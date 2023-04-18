package query

import (
	"fmt"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type GetUserRequest struct {
	ID *int
}

func GetUser(db db.Queryable, req GetUserRequest) (*structs.User, error) {
	// validate fields
	if req.ID == nil {
		return nil, fmt.Errorf("no id provided")
	}

	// use the list req
	users, err := ListUsers(db, ListUsersRequest{
		ID:     req.ID,
		Limit:  &[]int{1}[0],
		Offset: &[]int{0}[0],
	})
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users[0], nil
}
