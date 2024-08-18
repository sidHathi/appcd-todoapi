package models

type Attachment struct {
	Id string `json:"id"`
	TodoItemId string `json:"todo_item_id"`
	S3Url string `json:"s3_url"`
}