package services

import (
	"fmt"
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
)

func AddAttachment(itemId string, attachmentData models.CreateAttachmentData) (*models.Attachment, error) {
	var list_id string
	itemRow := db.Db.QueryRow("select list_id from todo_items where id=$1;", itemId)
	err := itemRow.Scan(&list_id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	atid := uuid.NewString()
	att := models.Attachment{
		Id:         atid,
		TodoItemId: itemId,
		ListId:     list_id,
		S3Url:      attachmentData.S3Url,
		FileType:   attachmentData.FileType,
	}
	_, err = db.Db.Exec("insert into attachments (id, list_id, item_id, s3_url, file_type) values ($1, $2, $3, $4, $5);", atid, list_id, itemId, attachmentData.S3Url, attachmentData.FileType)
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func GetAttachment(attachmentId string) (*models.Attachment, error) {
	var att models.Attachment
	row := db.Db.QueryRow("select item_id, list_id, s3_url, file_type from attachments where id=$1;", attachmentId)
	err := row.Scan(&att.TodoItemId, &att.ListId, &att.S3Url, &att.FileType)
	if err != nil {
		return nil, err
	}

	att.Id = attachmentId
	return &att, nil
}

func UpdateAttachment(attachmentId string, updateData models.CreateAttachmentData) error {
	currAtt, err := GetAttachment(attachmentId)
	if err != nil {
		return err
	}

	newUrl := currAtt.S3Url
	if updateData.S3Url != "" {
		newUrl = updateData.S3Url
	}
	newFileType := currAtt.FileType
	if updateData.FileType != "" {
		newUrl = updateData.FileType
	}
	_, err = db.Db.Exec("update attachments set s3_url=$1, file_type=$2 where id=$3", newUrl, newFileType, attachmentId)
	return err
}

func DeleteAttachment(id string) error {
	_, err := GetAttachment(id)
	if err != nil {
		return err
	}
	_, err = db.Db.Exec("delete from attachments where id=$1", id)
	return err
}
