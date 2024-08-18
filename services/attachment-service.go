package services

import (
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
)

func AddAttachment(itemId string, attachmentData models.CreateAttachmentData) (*models.Attachment, error) {
	var list_id string
	itemRow := db.Db.QueryRow("select (list_id) from todo_items where id=$1;", itemId)
	err := itemRow.Scan(&list_id)
	if err != nil {
		return nil, err
	}

	atid := uuid.NewString()
	att := models.Attachment {
		Id: atid,
		TodoItemId: itemId,
		ListId: list_id,
		S3Url: attachmentData.S3Url,
	}
	_, err = db.Db.Exec("insert into attachments (id, list_id, item_id, s3_url) values ($1, $2, $3, '$4');", atid, list_id, itemId, attachmentData.S3Url)
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func UpdateAttachment(attachmentId string, newUrl string) error {
	_, err := db.Db.Exec("update attachments set s3_url='$1' where id=$2", newUrl, attachmentId)
	return err
}

func DeleteAttachment(id string) error {
	_, err := db.Db.Exec("delete from attachments where id=$1", id)
	return err
}