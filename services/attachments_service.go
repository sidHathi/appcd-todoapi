package services

import (
	"todo-api/db"
	"todo-api/models"

	"github.com/google/uuid"
)

// Add an attachment to an item
func AddAttachment(itemId string, attachmentData models.CreateAttachmentData) (*models.Attachment, error) {
	// check to make sure the item exists
	var list_id string
	itemRow := db.Db.QueryRow("select list_id from todo_items where id=$1;", itemId)
	err := itemRow.Scan(&list_id)
	if err != nil {
		return nil, err
	}

	// create and add the attachment to the db
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

// get an attachment by its id
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

// update an attachment
func UpdateAttachment(attachmentId string, updateData models.CreateAttachmentData) error {
	// check to make sure it exists
	currAtt, err := GetAttachment(attachmentId)
	if err != nil {
		return err
	}

	// check to make sure the updateData is populated
	newUrl := currAtt.S3Url
	if updateData.S3Url != "" {
		newUrl = updateData.S3Url
	}
	newFileType := currAtt.FileType
	if updateData.FileType != "" {
		newUrl = updateData.FileType
	}
	// perform the sql query
	_, err = db.Db.Exec("update attachments set s3_url=$1, file_type=$2 where id=$3", newUrl, newFileType, attachmentId)
	return err
}

// delete an attachment
func DeleteAttachment(id string) error {
	_, err := GetAttachment(id)
	if err != nil {
		return err
	}
	_, err = db.Db.Exec("delete from attachments where id=$1", id)
	return err
}
