package query

import (
	"fmt"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type GetLinkRequest struct {
	ID         *int
	Identifier *string
}

func GetLink(db db.Queryable, req GetLinkRequest) (*structs.Link, error) {
	// validate the request
	if req.ID == nil && req.Identifier == nil {
		return nil, fmt.Errorf("no id or identifier provided")
	}

	// run list query to get the link
	links, err := ListLinks(db, ListLinksRequest{
		ID:     req.ID,
		Limit:  &[]int{1}[0],
		Offset: &[]int{0}[0],
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// check if the link exists
	if len(links) == 0 {
		return nil, nil
	}

	return links[0], nil

}
