package models

import (
	"fmt"
	"os"
	"time"

	"app.skyclerk.com/backend/library/files"
	"app.skyclerk.com/backend/library/store/object"
	"app.skyclerk.com/backend/services"
	"github.com/asaskevich/govalidator"
)

// File struct
type File struct {
	Id        uint      `gorm:"primary_key;column:FilesId" json:"id"`
	AccountId uint      `gorm:"column:FilesAccountId" sql:"not null" json:"account_id"`
	UpdatedAt time.Time `gorm:"column:FilesUpdatedAt" sql:"not null" json:"_"`
	CreatedAt time.Time `gorm:"column:FilesCreatedAt" sql:"not null" json:"_"`
	Host      string    `gorm:"column:FilesHost" sql:"not null" json:"_"`
	Name      string    `gorm:"column:FilesName" sql:"not null" json:"name"`
	Path      string    `gorm:"column:FilesPath" sql:"not null" json:"_"`
	ThumbPath string    `gorm:"column:FilesThumbPath" sql:"not null" json:"_"`
	Type      string    `gorm:"column:FilesType" sql:"not null" json:"type"`
	Hash      string    `gorm:"column:FilesHash" sql:"not null" json:"_"`
	Size      int64     `gorm:"column:FilesSize" sql:"not null" json:"size"`
	Url       string    `gorm:"-" json:"url"` // Not stored in DB.
}

//
// Set the table name.
//
func (File) TableName() string {
	return "Files"
}

//
// StoreFile - Store the file with our s3 file storage provider
//
func (t *DB) StoreFile(accountId uint, filePath string) (File, error) {
	// SafeFilename returns a cleaned-up filename that is safe to use.
	cleanedFileName := t.CleanFileName(filePath)

	// Get MD5 of the file.
	hash, err := files.Md5WithError(filePath)

	if err != nil {
		return File{}, err
	}

	// Get the file size.
	size, err := files.SizeWithError(filePath)

	if err != nil {
		return File{}, err
	}

	// Get the file type
	fileType, err := files.FileContentTypeWithError(filePath)

	if err != nil {
		return File{}, err
	}

	// Now that we have the file safely stored in our tmp directory time to process it.
	// First we create an entry in our files table so we know the ID.
	o := File{}
	o.Type = fileType
	o.Size = size
	o.Hash = hash
	o.Host = "amazon-s3"
	o.Name = cleanedFileName
	o.AccountId = accountId
	t.New().Save(&o)

	// Update the file path now that we have an id
	o.Path = fmt.Sprintf("accounts/%d/%d_%s", accountId, o.Id, cleanedFileName)
	t.New().Save(&o)

	// Upload file to our S3 store
	err = object.UploadObject(filePath, o.Path)

	if err != nil {
		return File{}, err
	}

	// TODO(spicer): Create thumbnail image.

	// Delete uploaded file
	err = os.Remove(filePath)

	if err != nil {
		services.Info(err)
	}

	// Return happy
	return o, nil
}

//
// GetSignedFileUrl - Pass in a path and get back a full url that is signed.
//
func (t *DB) GetSignedFileUrl(path string) string {
	//TODO(spicer): make this work

	return "https://app.skyclerk.com/" + path
}

//
// CleanFileName returns a cleaned-up filename that is safe to use.
//
func (t *DB) CleanFileName(fileName string) string {
	// This already checks for ";" so that ";;DROP TABLE x;" can't be part of the name.
	return govalidator.SafeFileName(fileName)
}

/* End File */