//
// Date: 2019-04-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"fmt"
	"go/build"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

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

	fmt.Println(w.Body.String())

	// // Grab result and convert to strut
	// result := models.File{}
	// err := json.Unmarshal([]byte(w.Body.String()), &result)
	//
	// // Query and get the ledger entry.
	// l, err := db.GetLedgerByAccountAndId(uint(33), ledger.Id)
	// st.Expect(t, err, nil)

	// // Test results
	// st.Expect(t, err, nil)
	// st.Expect(t, w.Code, 201)
	// st.Expect(t, result.Id, uint(1))
	// st.Expect(t, result.AccountId, uint(33))
	// st.Expect(t, result.Name, "boston-city-flow.jpg")
	// st.Expect(t, result.Type, "image/jpeg")
	// st.Expect(t, result.Size, int64(339773))
	// st.Expect(t, true, strings.Contains(result.Url, "https://cdn-dev.skyclerk.com/accounts/33/1_boston-city-flow.jpg?Expires="))
	// st.Expect(t, true, strings.Contains(result.Thumb600By600Url, "https://cdn-dev.skyclerk.com/accounts/33/1_thumb_600_600_boston-city-flow.jpg?Expires="))
	//
	// // Test ledger DB to file reults
	// st.Expect(t, l.Id, uint(1))
	// st.Expect(t, len(l.Files), 1)
	// st.Expect(t, l.Files[0].Name, "boston-city-flow.jpg")
	// st.Expect(t, l.Files[0].Size, int64(339773))
	// st.Expect(t, l.Files[0].Type, "image/jpeg")
	// st.Expect(t, l.Files[0].AccountId, uint(33))
	// st.Expect(t, true, strings.Contains(l.Files[0].Url, "https://cdn-dev.skyclerk.com/accounts/33/1_boston-city-flow.jpg?Expires="))
	// st.Expect(t, true, strings.Contains(l.Files[0].Thumb600By600Url, "https://cdn-dev.skyclerk.com/accounts/33/1_thumb_600_600_boston-city-flow.jpg?Expires="))
	//
	// // If we are testing locally (not on CI) we test to see if the file is on AWS with our signed key
	// if len(os.Getenv("AWS_CLOUDFRONT_PRIVATE_SIGN_KEY")) > 0 {
	// 	// Make sure the file is not already there.
	// 	os.Remove("/tmp/1_boston-city-flow.jpg")
	//
	// 	err := downloadFile("/tmp/1_boston-city-flow.jpg", result.Url)
	// 	st.Expect(t, err, nil)
	//
	// 	// Get the MD5 hash from the DB and compare
	// 	ff := models.File{}
	// 	db.New().Find(&ff, result.Id)
	// 	st.Expect(t, ff.Hash, files.Md5("/tmp/1_boston-city-flow.jpg"))
	//
	// 	err = os.Remove("/tmp/1_boston-city-flow.jpg")
	// 	st.Expect(t, err, nil)
	//
	// 	// --- here we chnage the key. We are just verifying AWS is honoring the singing
	//
	// 	err = downloadFile("/tmp/1_boston-city-flow.jpg", result.Url+"blah")
	// 	st.Expect(t, err, nil)
	//
	// 	// Hash should be of "unauthorized"
	// 	st.Expect(t, "b741482e554b40fd711a012fa74461cd", files.Md5("/tmp/1_boston-city-flow.jpg"))
	//
	// 	err = os.Remove("/tmp/1_boston-city-flow.jpg")
	// 	st.Expect(t, err, nil)
	//
	// 	// ---- Test the thumb nail
	//
	// 	// Make sure the file is not already there.
	// 	os.Remove("/tmp/1_thumb_800_800_boston-city-flow.jpg")
	//
	// 	err = downloadFile("/tmp/1_thumb_600_600_boston-city-flow.jpg", result.Thumb600By600Url)
	// 	st.Expect(t, err, nil)
	// 	st.Expect(t, "ef6363908d9bdd27c9b1737ad26afbc6", files.Md5("/tmp/1_thumb_600_600_boston-city-flow.jpg"))
	//
	// 	err = os.Remove("/tmp/1_thumb_600_600_boston-city-flow.jpg")
	// 	st.Expect(t, err, nil)
	// }
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
