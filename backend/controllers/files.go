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

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/files"
	"app.skyclerk.com/backend/library/store/object"
	"app.skyclerk.com/backend/models"
)

//
// CreateFile - Upload a file to the account.
//
func (t *Controller) CreateFile(c *gin.Context) {
	// Options fields that can be included in the post for later assignment.
	id := c.PostForm("id")
	table := c.PostForm("object")

	// Get user id.
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
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	// Save the uploaded file. Store file in tmp directory
	filePath := fmt.Sprintf("%s/%s", cacheDir, filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	// SafeFilename returns a cleaned-up filename that is safe to use.
	cleanedFileName := t.db.CleanFileName(file.Filename)

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields name=%s and email=%s.", file.Filename, table, id))

	// Get MD5 of the file.
	hash, err := files.Md5WithError(filePath)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	// Get the file size.
	size, err := files.SizeWithError(filePath)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	// Get the file type
	fileType, err := files.FileContentTypeWithError(filePath)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	// Now that we have the file safely stored in our tmp directory time to process it.
	// First we create an entry in our files table so we know the ID.
	fileModel := models.File{}
	fileModel.Type = fileType
	fileModel.Size = size
	fileModel.Hash = hash
	fileModel.Host = "amazon-s3"
	fileModel.Name = cleanedFileName
	fileModel.AccountId = accountId
	t.db.New().Save(&fileModel)

	// Update the file path now that we have an id
	fileModel.Path = fmt.Sprintf("accounts/%d/%d_%s", accountId, fileModel.Id, cleanedFileName)
	t.db.New().Save(&fileModel)

	// Upload file to our S3 store
	err = object.UploadObject(filePath, fileModel.Path)

	if err != nil {
		panic(err)
	}

	// TODO(spicer): Create thumbnail image.

	// // Setup Label obj
	// o := models.Label{}
	//
	// // Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	// if t.ValidateRequest(c, &o, "create") != nil {
	// 	return
	// }
	//
	// // Make sure the AccountId is correct.
	// o.AccountId = uint(c.MustGet("accountId").(int))
	//
	// // Clean up some vars
	// o.Name = strings.Trim(o.Name, " ")
	//
	// // Create label
	// t.db.New().Create(&o)
	//
	// // Return happy.
	// response.RespondCreated(c, o, nil)
}

/* End File */
