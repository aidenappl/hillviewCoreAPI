package structs

type Link struct {
	ID          int64  `json:"id"`
	Route       string `json:"route"`
	Destination string `json:"destination"`
	Active      bool   `json:"active"`
	Creator     UserTS `json:"creator"`
	Clicks      int    `json:"clicks"`
	InsertedAt  string `json:"inserted_at"`
}
