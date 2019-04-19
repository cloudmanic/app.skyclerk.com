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
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"

	"github.com/cloudmanic/app.skyclerk.com/backend/library/store/object"
)

//
// CreateFile - Upload a file to the account.
//
func (t *Controller) CreateFile(c *gin.Context) {
	// Options fields that can be included in the post for later assignment.
	id := c.PostForm("id")
	table := c.PostForm("object")

	// This is the file we are uploading.
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	// Save the uploaded file. Store file in tmp directory
	filename := "/tmp/" + filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully with fields name=%s and email=%s.", file.Filename, table, id))

	// Now that we have the file safely stored in our tmp directory time to process it.

	list, err := object.ListObjects("")

	if err != nil {
		panic(err)
	}

	spew.Dump(list)

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
