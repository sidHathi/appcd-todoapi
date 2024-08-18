package models

type Attachment struct {
	Id string `json:"id"`
	TodoItemId string `json:"todo_item_id"`
	ListId string `json:"list_id"`
	S3Url string `json:"s3_url"`
}

type CreateAttachmentData struct {
	S3Url string `json:"s3_url"`
}