package models

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	"github.com/disintegration/imaging"

	"app.skyclerk.com/backend/library/files"
	"app.skyclerk.com/backend/library/store/object"
	"app.skyclerk.com/backend/services"
)

// File struct
type File struct {
	Id               uint      `gorm:"primary_key;column:FilesId" json:"id"`
	AccountId        uint      `gorm:"column:FilesAccountId" sql:"not null" json:"account_id"`
	UpdatedAt        time.Time `gorm:"column:FilesUpdatedAt" sql:"not null" json:"_"`
	CreatedAt        time.Time `gorm:"column:FilesCreatedAt" sql:"not null" json:"_"`
	Host             string    `gorm:"column:FilesHost" sql:"not null" json:"_"`
	Name             string    `gorm:"column:FilesName" sql:"not null" json:"name"`
	Path             string    `gorm:"column:FilesPath" sql:"not null" json:"_"`
	ThumbPath        string    `gorm:"column:FilesThumbPath" sql:"not null" json:"_"` // TODO(spicer): rename this Thumb800x800
	Type             string    `gorm:"column:FilesType" sql:"not null" json:"type"`
	Hash             string    `gorm:"column:FilesHash" sql:"not null" json:"_"`
	Size             int64     `gorm:"column:FilesSize" sql:"not null" json:"size"`
	Url              string    `gorm:"-" json:"url"`                  // Not stored in DB.
	Thumb800By800Url string    `gorm:"-" json:"thumb_800_by_800_url"` // Not stored in DB.
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
		services.Info(err)
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

	// Set upload path
	up := fmt.Sprintf("accounts/%d/%d_%s", accountId, o.Id, cleanedFileName)

	// Upload file to our S3 store
	err = object.UploadObject(filePath, up)

	if err != nil {
		services.Critical(errors.New(fmt.Sprintf("FileId: %d, AccountId: %d Error: %s", o.Id, accountId, err.Error())))
		return File{}, err
	}

	// Update the file path now that we have an id
	o.Path = up
	t.New().Save(&o)

	// Create and store at S3 the thumbnail image.
	err = t.CreateAndStoreThumbnailImage(&o, cleanedFileName, filePath)

	if err != nil {
		services.Critical(err)
	}

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
// This url is good for 5 mins.
//
func (t *DB) GetSignedFileUrl(path string) string {
	// RUL we need to sign.
	rawURL := os.Getenv("OBJECT_BASE_URL") + "/" + path

	// This is a small hack for testing. This is because we do not want to share our keys with CI
	if len(os.Getenv("AWS_CLOUDFRONT_PRIVATE_SIGN_KEY")) == 0 {
		return rawURL + "?Expires="
	}

	// Decode the base64 and pass in a real private key
	sDec, _ := base64.StdEncoding.DecodeString(os.Getenv("AWS_CLOUDFRONT_PRIVATE_SIGN_KEY"))

	// Build the private key obj
	block, _ := pem.Decode(sDec)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	// Sign URL to be valid for 5 mins from now.
	signer := sign.NewURLSigner(os.Getenv("AWS_CLOUDFRONT_KEY_ID"), key)
	signedURL, err := signer.Sign(rawURL, time.Now().Add(5*time.Minute))

	if err != nil {
		services.Info(err)
	}

	// Return happy
	return signedURL
}

//
// CleanFileName returns a cleaned-up filename that is safe to use.
//
func (t *DB) CleanFileName(fileName string) string {
	// This already checks for ";" so that ";;DROP TABLE x;" can't be part of the name.
	return govalidator.SafeFileName(fileName)
}

//
// Create thumbnail image
//
func (t *DB) CreateAndStoreThumbnailImage(file *File, cleanedFileName string, filePath string) error {
	// Thumb base name
	tbn := fmt.Sprintf("%d_thumb_800_800_%s", file.Id, cleanedFileName)

	// File cache dir.
	cacheDir := fmt.Sprintf("%s/thumbs/%d", os.Getenv("CACHE_DIR"), file.AccountId)

	// Make the directory we store this file to
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}

	// TODO(spicer): Create thumbnail of PDF types

	// Open a test image.
	src, err := imaging.Open(filePath)
	if err != nil {
		return err
	}

	// Resize and crop the srcImage to fill the 100x100px area.
	dst := imaging.Fill(src, 800, 800, imaging.Center, imaging.Lanczos)

	// Saved thumbe file path
	tbfp := cacheDir + "/" + tbn

	// Save the resulting image as JPEG.
	err = imaging.Save(dst, tbfp)
	if err != nil {
		return err
	}

	// Set thumb path
	tp := fmt.Sprintf("accounts/%d/%d_thumb_800_800_%s", file.AccountId, file.Id, cleanedFileName)

	// Upload file to S3
	err = object.UploadObject(tbfp, tp)

	if err != nil {
		services.Critical(errors.New(fmt.Sprintf("Thumbnail FileId: %d, AccountId: %d Error: %s", file.Id, file.AccountId, err.Error())))
		return err
	}

	// Update the file path now that we have an id
	file.ThumbPath = tp
	t.New().Save(&file)

	// Delete thumb file
	err = os.Remove(tbfp)

	if err != nil {
		services.Info(err)
	}

	// Return  happy
	return nil
}

/* End File */
