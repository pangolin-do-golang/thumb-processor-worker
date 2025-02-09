package domain

type Event struct {
	ID        string `json:"id"`
	VideoPath string `json:"video"`
	ThumbPath string `json:"thumbnail"`
	Metadata  interface{}
}
