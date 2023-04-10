package structs

var (
	VideoStatusDeleted = 4
)

type Video struct {
	ID          int         `json:"id"`
	UUID        string      `json:"uuid"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Thumbnail   string      `json:"thumbnail"`
	URL         string      `json:"url"`
	Status      *GeneralNSN `json:"status"`
	InsertedAt  string      `json:"inserted_at"`
}
