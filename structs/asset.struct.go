package structs

import "time"

type Asset struct {
	ID          int64                 `json:"id"`
	Name        *string               `json:"name"`
	ImageURL    *string               `json:"image_url"`
	Identifier  *string               `json:"identifier"`
	Description *string               `json:"description"`
	Category    *GeneralNSN           `json:"category"`
	Status      *GeneralNSN           `json:"status"`
	Metadata    *AssetMetadata        `json:"metadata"`
	ActiveTab   *AssetCheckoutOmitted `json:"active_tab"`
	InsertedAt  time.Time             `json:"inserted_at"`
}
