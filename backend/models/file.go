package models

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	Assigned         int       `gorm:"column:FilesAssigned" sql:"not null" json:"_"`
	Url              string    `gorm:"-" json:"url"`                  // Not stored in DB.
	Thumb600By600Url string    `gorm:"-" json:"thumb_600_by_600_url"` // Not stored in DB.
}

//
// Set the table name.
//
func (File) TableName() string {
	return "Files"
}

//
// GetFileByAccountAndId by account and id.
//
func (db *DB) GetFileByAccountAndId(accountId uint, id uint) (File, error) {
	// File to return
	c := File{}

	// Make query
	if db.New().Where("FilesAccountId = ? AND FilesId = ?", accountId, id).First(&c).RecordNotFound() {
		return File{}, errors.New("File entry not found.")
	}

	// Add in a signed URL
	c.Url = db.GetSignedFileUrl(c.Path)
	c.Thumb600By600Url = db.GetSignedFileUrl(c.ThumbPath)

	// Return result
	return c, nil
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
	o.Assigned = 1 // TODO(spicer): kill this after we kill PHP
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
	err = t.CreateAndStoreThumbnailImage(&o, cleanedFileName, filePath, fileType)

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
	// TODO(spicer): If we do not have a cloudfront key we revert to using S3 bucket signing (good for testing)
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
func (t *DB) CreateAndStoreThumbnailImage(file *File, cleanedFileName string, filePath string, fileType string) error {
	const width int = 600
	const height int = 600

	// File cache dir.
	cacheDir := fmt.Sprintf("%s/thumbs/%d", os.Getenv("CACHE_DIR"), file.AccountId)

	// Make the directory we store this file to
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}

	// The path to the thumb nail we are going to upload.
	var tbfp string

	// If this is a PDF we use imaginary.skyclerk.com to create the thumbnail - This is hacky to support testing.
	if fileType == "application/pdf" && (len(os.Getenv("AWS_CLOUDFRONT_PRIVATE_SIGN_KEY")) > 0) {
		t, err := t.GetPdfThumbNail(file, width, height, cleanedFileName)

		if err != nil {
			return err
		}

		tbfp = t
	}

	// If this is an image we create a thumbnail a different way.
	if (fileType == "image/jpeg") || (fileType == "image/png") || (fileType == "image/gif") {
		t, err2 := t.GetImageThumbNail(file, filePath, width, height, cleanedFileName)

		if err2 != nil {
			return err2
		}

		tbfp = t
	}

	// If we do not have a thumbnail no need to go on.
	if len(tbfp) == 0 {
		err := errors.New(fmt.Sprintf("Thumbnail Failed to create FileId: %d, AccountId: %d Error: %s", file.Id, file.AccountId, "Unable to create thumbnail."))
		return err
	}

	// Set thumb path
	tp := fmt.Sprintf("accounts/%d/%s", file.AccountId, filepath.Base(tbfp))

	// Upload file to S3
	err := object.UploadObject(tbfp, tp)

	if err != nil {
		return errors.New(fmt.Sprintf("Thumbnail FileId: %d, AccountId: %d Error: %s", file.Id, file.AccountId, err.Error()))
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

//
// GetImageThumbNail get a thumbnail from an image.
//
func (t *DB) GetImageThumbNail(file *File, filePath string, width int, height int, cleanedFileName string) (string, error) {
	// File cache dir.
	cacheDir := fmt.Sprintf("%s/thumbs/%d", os.Getenv("CACHE_DIR"), file.AccountId)

	// Thumb base name
	cleanedFileName = strings.Replace(cleanedFileName, "pdf", "jpeg", 100)
	cleanedFileName = strings.Replace(cleanedFileName, "PDF", "jpeg", 100)
	tbn := fmt.Sprintf("%d_thumb_%d_%d_%s", file.Id, width, height, cleanedFileName)

	// Open a test image.
	src, err := imaging.Open(filePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("GetImageThumbNail Failed (#001) - FileId: %d, AccountId: %d Error: %s", file.Id, file.AccountId, err.Error()))
	}

	// Resize and crop the srcImage to fill the widthXheight area.
	dst := imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)

	// Saved thumbe file path
	tbfp := cacheDir + "/" + tbn

	// Save the resulting image as JPEG.
	err = imaging.Save(dst, tbfp)
	if err != nil {
		return "", errors.New(fmt.Sprintf("GetImageThumbNail Failed (#002) - FileId: %d, AccountId: %d Error: %s", file.Id, file.AccountId, err.Error()))
	}

	// Return happy
	return tbfp, nil
}

//
// GetPdfThumbNail uses imaginary.skyclerk.com to build a thumbnail of a pdf.
//
func (t *DB) GetPdfThumbNail(file *File, width int, height int, cleanedFileName string) (string, error) {
	// Thumb base name
	cleanedFileName = strings.Replace(cleanedFileName, "pdf", "jpeg", 100)
	cleanedFileName = strings.Replace(cleanedFileName, "PDF", "jpeg", 100)
	tbn := fmt.Sprintf("%d_thumb_%d_%d_%s", file.Id, width, height, cleanedFileName)

	// File cache dir.
	cacheDir := fmt.Sprintf("%s/thumbs/%d", os.Getenv("CACHE_DIR"), file.AccountId)

	// Set a signed url of the orginal
	signedUrl := t.GetSignedFileUrl(file.Path)

	request := fmt.Sprintf("%s/convert?type=jpeg&width=%d&height=%d&quality=95&gravity=smart&url=%s", os.Getenv("IMAGINARY_HOST"), width, height, url.QueryEscape(signedUrl))

	// Build imaginary request
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", request, nil)

	// Headers
	req.Header.Add("API-Key", os.Getenv("IMAGINARY_KEY"))

	err = req.ParseForm()
	if err != nil {
		return "", errors.New(fmt.Sprintf("getPdfThumbNail Failed (#001) - FileId: %d, AccountId: %d Error: %s", file.Id, file.AccountId, err.Error()))
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		return "", errors.New(fmt.Sprintf("getPdfThumbNail Failed (#002) - FileId: %d, AccountId: %d Error: %s", file.Id, file.AccountId, err.Error()))
	}

	// Must have a status code of 200 or something failed.
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("getPdfThumbNail Failed (#003) - FileId: %d, AccountId: %d Status Code: %d", file.Id, file.AccountId, resp.StatusCode))
	}

	// Saved thumbe file path
	tbfp := cacheDir + "/" + tbn

	// Create the file to store locally
	out, err := os.Create(tbfp)
	if err != nil {
		return "", errors.New(fmt.Sprintf("getPdfThumbNail Failed (#004) - FileId: %d, AccountId: %d Error: %s", file.Id, file.AccountId, err.Error()))
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return "", errors.New(fmt.Sprintf("getPdfThumbNail Failed (#005) - FileId: %d, AccountId: %d Error: %s", file.Id, file.AccountId, err.Error()))
	}

	// Return happy
	return tbfp, nil
}

/* End File */
