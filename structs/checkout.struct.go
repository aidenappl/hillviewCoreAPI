package structs

import "time"

type Checkout struct {
	ID             int         `json:"id"`
	User           *MobileUser `json:"user"`
	AssociatedUser *int        `json:"associated_user,omitempty"`
	Asset          *Asset      `json:"asset"`
	AssetID        *int        `json:"asset_id,omitempty"`
	Offsite        int         `json:"offsite"`
	CheckoutStatus *GeneralNSN `json:"checkout_status"`
	CheckoutNotes  *string     `json:"checkout_notes"`
	TimeOut        *time.Time  `json:"time_out"`
	TimeIn         *time.Time  `json:"time_in"`
	ExpectedIn     *time.Time  `json:"expected_in"`
}
