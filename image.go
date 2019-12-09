package spree

import "time"

type ImageId int
type Image struct {
	Id                    ImageId   `json:"id"`
	Position              int       `json:"position"`
	AttachmentContentType string    `json:"attachment_content_type"`
	AttachmentFileName    string    `json:"attachment_file_name"`
	Type                  string    `json:"type"`
	AttachmentUpdatedAt   time.Time `json:"attachment_updated_at"`
	AttachmentWidth       int       `json:"attachment_width"`
	AttachmentHeight      int       `json:"attachment_height"`
	Alt                   string    `json:"alt"`
	ViewableType          string    `json:"viewable_type"`
	ViewableId            int       `json:"viewable_id"`
	MiniURL               string    `json:"mini_url"`
	SmallURL              string    `json:"small_url"`
	ProductURL            string    `json:"product_url"`
	LargeURL              string    `json:"large_url"`
}
