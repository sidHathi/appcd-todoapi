package models

// Attachment represents a file attachment in the system
// we expect attachments to be stored using s3 urls to
// where the data is stored in AWS
// swagger:model
type Attachment struct {
	Id         string `json:"id"`
	TodoItemId string `json:"todo_item_id"`
	ListId     string `json:"list_id"`
	S3Url      string `json:"s3_url"`
	FileType   string `json:"file_type"`
}

// CreateAttachmentData represents the data needed to create
// a new attachment
// swagger:model
type CreateAttachmentData struct {
	S3Url    string `json:"s3_url"`
	FileType string `json:"file_type"`
}
