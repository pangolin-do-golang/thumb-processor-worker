package domain

type Event struct {
	ID        string `json:"id"`
	UserEmail string `json:"user_email"`
	VideoPath string `json:"video"`
	ThumbPath string `json:"thumbnail"`
	Metadata  interface{}
}
