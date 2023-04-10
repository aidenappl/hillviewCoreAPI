package structs

type AssetCheckoutOmitted struct {
	ID             int `json:"id"`
	AssociatedUser int `json:"associated_user"`
}

type AssetMetadata struct {
	SerialNumber *string `json:"serial_number"`
	Manufacturer *string `json:"manufacturer"`
	Model        *string `json:"model"`
	Notes        *string `json:"notes"`
}

type GeneralNSN struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}
