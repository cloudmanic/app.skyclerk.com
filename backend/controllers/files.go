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

	"app.skyclerk.com/backend/library/files"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/services"
	"github.com/gin-gonic/gin"
)

// MaxFileUploadSize in bytes
const maxFileUploadSize int64 = 50000000 // 50 Megabytes

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

	// TODO(spicer): Validate max file size.
	size, err := files.SizeWithError(filePath)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "An error happend when uploading file (#002). Please contact help@skyclerk.com."})
		return
	}

	if size > maxFileUploadSize {
		// TODO(spicer): Validate max file size.
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("We have a %d megabyte upload limit.", maxFileUploadSize)})
	}

	// Store the file with S3 and create Files entry.
	o, err := t.db.StoreFile(accountId, filePath)
	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "An error happend when uploading file (#003). Please contact help@skyclerk.com."})
		return
	}

	// Add in a signed URL
	o.Url = t.db.GetSignedFileUrl(o.Path)

	// Return happy.
	response.RespondCreated(c, o, nil)
}

/* End File */
