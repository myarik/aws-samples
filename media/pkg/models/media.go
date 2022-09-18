package models

type MediaType string

const (
	Image MediaType = "image"
	Video MediaType = "video"
)

type Media struct {
	Id        string    `json:"id"`
	ProductId string    `json:"product_id"`
	Type      MediaType `json:"type"`
	URL       string    `json:"url"`
	CreatedAt int64     `json:"timestamp"`
}

type Thumbnail struct {
	ProductId string `json:"product_id"`
	MediaId   string `json:"media_id"`
	URL       string `json:"url"`
	Height    int    `json:"height"`
	Width     int    `json:"width"`
}
