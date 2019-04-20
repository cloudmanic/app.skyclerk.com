//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-29
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/adelowo/filer"
	"github.com/adelowo/filer/validator"
	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

//
// CreateFile - Upload a file to the account.
//
func (t *Controller) CreateFile(c *gin.Context) {
	// Options fields that can be included in the post for later assignment.
	// id := c.PostForm("id")
	// table := c.PostForm("object")

	// Get user id. f
	//userId := uint(c.MustGet("userId").(int))

	// AccountId.
	accountId := uint(c.MustGet("accountId").(int))

	// File cache dir.
	cacheDir := fmt.Sprintf("%s/uploads/%d", os.Getenv("CACHE_DIR"), accountId)

	// Make the directory we store this file to
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}

	// This is the file we are uploading.
	file, err := c.FormFile("file")
	if err != nil {
		// TODO(spicer): Validate max file size.
		c.JSON(http.StatusBadRequest, gin.H{"error": "A file is required."})
		return
	}

	// Save the uploaded file. Store file in tmp directory
	filePath := fmt.Sprintf("%s/%s", cacheDir, filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "An error happend when uploading file (#001). Please contact help@skyclerk.com."})
		return
	}

	// Validate file we are uploading. JSON error set in function.
	o, err := t.validateUploadedFile(c, filePath, accountId)
	if err != nil {
		return
	}

	// Add in a signed URL
	o.Url = t.db.GetSignedFileUrl(o.Path)

	// Return happy.
	response.RespondCreated(c, o, nil)
}

//
// validateUploadedFile - Validate a file we are uploading.
//
func (t *Controller) validateUploadedFile(c *gin.Context, filePath string, accountId uint) (models.File, error) {
	// Setup validators
	max, _ := filer.LengthInBytes("50MB")
	min, _ := filer.LengthInBytes("1B")
	val := validator.NewSizeValidator(max, min)
	val2 := validator.NewMimeTypeValidator([]string{"image/jpeg", "image/png", "image/gif", "application/pdf"})

	// Open file so we can validate
	vf, _ := os.Open(filePath)

	// Validate max size
	if _, err := val.Validate(vf); err != nil {
		// TODO(spicer): Validate max file size.
		c.JSON(http.StatusBadRequest, gin.H{"error": "We have a 50MB upload limit."})
		return models.File{}, err
	}

	// Validate file type
	if _, err := val2.Validate(vf); err != nil {
		// TODO(spicer): Validate max file size.
		c.JSON(http.StatusBadRequest, gin.H{"error": "We only allow image and pdf files to be uploaded."})
		return models.File{}, err
	}

	// Store the file with S3 and create Files entry.
	o, err := t.db.StoreFile(accountId, filePath)
	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "An error happend when uploading file (#003). Please contact help@skyclerk.com."})
		return models.File{}, err
	}

	// Return happy
	return o, nil
}

/* End File */
