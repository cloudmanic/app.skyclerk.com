//
// Date: 2019-04-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestCreateSnapClerk01 - Test create snapclerk 01
//
func TestCreateSnapClerk01(t *testing.T) {
	// test file.
	testFile := build.Default.GOPATH + "/src/app.skyclerk.com/backend/library/test/files/Image 2019-04-19 at 10.10.22 AM.png"

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Add in test snapclerks so all ids are not "1"
	s1 := test.GetRandomSnapClerk(44)
	db.New().Save(&s1)

	s2 := test.GetRandomSnapClerk(55)
	db.New().Save(&s2)

	s3 := test.GetRandomSnapClerk(44)
	db.New().Save(&s3)

	// Build file to post
	buffer, writer := buildSnapClerkFileForm(t, testFile, 55.23, "ABC, Inc.", "Marketing", "label #1, label #2, label #3", "1233.2342", "4234.3242", "This is a test note for snapclerk.")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/44/snapclerk", buffer)

	// Set content type header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 44)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/:account/snapclerk", c.CreateSnapClerk)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.SnapClerk{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)
	st.Expect(t, err, nil)

	// Query and get the SnapClerk entry.
	l, err := db.GetSnapClerkByAccountAndId(result.AccountId, result.Id)
	st.Expect(t, err, nil)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Id, uint(4))
	st.Expect(t, result.AccountId, uint(44))
	st.Expect(t, result.Status, "Pending")
	st.Expect(t, result.Contact, "ABC, Inc.")
	st.Expect(t, result.Category, "Marketing")
	st.Expect(t, result.Labels, "label #1, label #2, label #3")
	st.Expect(t, result.Note, "This is a test note for snapclerk.")
	st.Expect(t, result.Lat, "1233.2342")
	st.Expect(t, result.Lon, "4234.3242")
	st.Expect(t, result.File.Id, uint(1))
	st.Expect(t, result.File.Name, "image-2019-04-19-at-10.10.22-am.png")
	st.Expect(t, result.File.Type, "image/png")
	st.Expect(t, result.File.Size, int64(861591))
	st.Expect(t, true, strings.Contains(result.File.Url, "https://cdn-dev.skyclerk.com/accounts/44/1_image-2019-04-19-at-10.10.22-am.png?Expires="))
	st.Expect(t, true, strings.Contains(result.File.Thumb600By600Url, "https://cdn-dev.skyclerk.com/accounts/44/1_thumb_600_600_image-2019-04-19-at-10.10.22-am.png?Expires="))

	// Test SnapClerk DB to file reults
	st.Expect(t, l.Id, uint(4))
	st.Expect(t, l.AccountId, uint(44))
	st.Expect(t, l.AddedById, uint(109))
	st.Expect(t, l.Status, "Pending")
	st.Expect(t, l.Contact, "ABC, Inc.")
	st.Expect(t, l.Category, "Marketing")
	st.Expect(t, l.Labels, "label #1, label #2, label #3")
	st.Expect(t, l.Note, "This is a test note for snapclerk.")
	st.Expect(t, l.Lat, "1233.2342")
	st.Expect(t, l.Lon, "4234.3242")
	st.Expect(t, l.File.Id, uint(1))
	st.Expect(t, l.File.Name, "image-2019-04-19-at-10.10.22-am.png")
	st.Expect(t, l.File.Type, "image/png")
	st.Expect(t, l.File.Size, int64(861591))
	st.Expect(t, true, strings.Contains(l.File.Url, "https://cdn-dev.skyclerk.com/accounts/44/1_image-2019-04-19-at-10.10.22-am.png?Expires="))
	st.Expect(t, true, strings.Contains(l.File.Thumb600By600Url, "https://cdn-dev.skyclerk.com/accounts/44/1_thumb_600_600_image-2019-04-19-at-10.10.22-am.png?Expires="))
}

//
// buildSnapClerkFileForm so we can post a file.
//
func buildSnapClerkFileForm(t *testing.T, filePath string, amount float64, contact string, category string, labels string, lat string, lon string, note string) (*bytes.Buffer, *multipart.Writer) {
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

	// Build snapclerk fields
	err = writer.WriteField("amount", fmt.Sprintf("%f", amount))
	st.Expect(t, err, nil)

	err = writer.WriteField("contact", fmt.Sprintf("%s", contact))
	st.Expect(t, err, nil)

	err = writer.WriteField("category", fmt.Sprintf("%s", category))
	st.Expect(t, err, nil)

	err = writer.WriteField("labels", fmt.Sprintf("%s", labels))
	st.Expect(t, err, nil)

	err = writer.WriteField("lat", fmt.Sprintf("%s", lat))
	st.Expect(t, err, nil)

	err = writer.WriteField("lon", fmt.Sprintf("%s", lon))
	st.Expect(t, err, nil)

	err = writer.WriteField("note", fmt.Sprintf("%s", note))
	st.Expect(t, err, nil)

	// Close writer
	err = writer.Close()
	st.Expect(t, err, nil)

	return body, writer
}

/* End File */
