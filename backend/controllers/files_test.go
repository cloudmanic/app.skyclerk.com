//
// Date: 2018-03-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-29
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"go/build"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nbio/st"

	"app.skyclerk.com/backend/models"
)

//
// Test create File 01
//
func TestCreateFiles01(t *testing.T) {
	// test file.
	testFile := build.Default.GOPATH + "/src/app.skyclerk.com/backend/library/test/files/Boston City Flow.jpg"

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Build file to post
	buffer, writer := buildLedgerFileform(t, testFile, "ledger", uint(88))

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/files", buffer)

	// Set content type header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/:account/files", c.CreateFile)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.File{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Id, uint(1))
	st.Expect(t, result.AccountId, uint(33))
	st.Expect(t, result.Name, "boston-city-flow.jpg")
	st.Expect(t, result.Type, "image/jpeg")
	st.Expect(t, result.Size, int64(339773))
	st.Expect(t, result.Url, "https://app.skyclerk.com/accounts/33/1_boston-city-flow.jpg")
}

//
// buildLedgerFileform so we can pust a file.
//
func buildLedgerFileform(t *testing.T, filePath string, object string, id uint) (*bytes.Buffer, *multipart.Writer) {
	// Build buffer for file to upload.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create form file body
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	st.Expect(t, err, nil)

	// Open file handle
	fh, err := os.Open(filePath)
	st.Expect(t, err, nil)
	defer fh.Close()

	// Copy file data to form body.
	_, err = io.Copy(part, fh)
	st.Expect(t, err, nil)

	// Add in a field to tell which object this should be attached to.
	err = writer.WriteField("object", object)
	st.Expect(t, err, nil)

	// Add in a field to tell the id of the object
	err = writer.WriteField("id", strconv.Itoa(int(id)))
	st.Expect(t, err, nil)

	// Close writer
	err = writer.Close()
	st.Expect(t, err, nil)

	return body, writer
}

/* End File */
