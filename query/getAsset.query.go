package query

import (
	"fmt"

	"github.com/hillview.tv/coreAPI/db"
	"github.com/hillview.tv/coreAPI/structs"
)

type GetAssetRequest struct {
	ID         *int
	Identifier *string
}

func GetAsset(db db.Queryable, req GetAssetRequest) (*structs.Asset, error) {
	// validate the request
	if req.ID == nil && req.Identifier == nil {
		return nil, fmt.Errorf("no id or identifier provided")
	}

	// run list query to get the asset
	assets, err := ListAssets(db, ListAssetsRequest{
		ID:         req.ID,
		Identifier: req.Identifier,
		UseOr:      true,
		Limit:      1,
		Offset:     0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// check if the asset exists
	if len(*assets) == 0 {
		return nil, nil
	}

	return &(*assets)[0], nil

}
