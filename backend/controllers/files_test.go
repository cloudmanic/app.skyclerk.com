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
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nbio/st"

	"app.skyclerk.com/backend/library/files"
	"app.skyclerk.com/backend/models"
)

//
// Test create File 01
//
func TestCreateFiles01(t *testing.T) {
	// test file.
	testFile := build.Default.GOPATH + "/src/app.skyclerk.com/backend/library/test/files/Boston City Flow.jpg"

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Build file to post
	buffer, writer := buildLedgerFileform(t, testFile)

	// Attach a ledger  to add this file to.
	err := writer.WriteField("ledger_id", "55")
	st.Expect(t, err, nil)

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
	err = json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Id, uint(1))
	st.Expect(t, result.AccountId, uint(33))
	st.Expect(t, result.Name, "boston-city-flow.jpg")
	st.Expect(t, result.Type, "image/jpeg")
	st.Expect(t, result.Size, int64(339773))
	st.Expect(t, true, strings.Contains(result.Url, "https://cdn-dev.skyclerk.com/accounts/33/1_boston-city-flow.jpg?Expires="))
	st.Expect(t, true, strings.Contains(result.Thumb600By600Url, "https://cdn-dev.skyclerk.com/accounts/33/1_thumb_600_600_boston-city-flow.jpg?Expires="))

	// If we are testing locally (not on CI) we test to see if the file is on AWS with our signed key
	if len(os.Getenv("AWS_CLOUDFRONT_PRIVATE_SIGN_KEY")) > 0 {
		// Make sure the file is not already there.
		os.Remove("/tmp/1_boston-city-flow.jpg")

		err := downloadFile("/tmp/1_boston-city-flow.jpg", result.Url)
		st.Expect(t, err, nil)

		// Get the MD5 hash from the DB and compare
		ff := models.File{}
		db.New().Find(&ff, result.Id)
		st.Expect(t, ff.Hash, files.Md5("/tmp/1_boston-city-flow.jpg"))

		err = os.Remove("/tmp/1_boston-city-flow.jpg")
		st.Expect(t, err, nil)

		// --- here we chnage the key. We are just verifying AWS is honoring the singing

		err = downloadFile("/tmp/1_boston-city-flow.jpg", result.Url+"blah")
		st.Expect(t, err, nil)

		// Hash should be of "unauthorized"
		st.Expect(t, "b741482e554b40fd711a012fa74461cd", files.Md5("/tmp/1_boston-city-flow.jpg"))

		err = os.Remove("/tmp/1_boston-city-flow.jpg")
		st.Expect(t, err, nil)

		// ---- Test the thumb nail

		// Make sure the file is not already there.
		os.Remove("/tmp/1_thumb_800_800_boston-city-flow.jpg")

		err = downloadFile("/tmp/1_thumb_600_600_boston-city-flow.jpg", result.Thumb600By600Url)
		st.Expect(t, err, nil)
		st.Expect(t, "ef6363908d9bdd27c9b1737ad26afbc6", files.Md5("/tmp/1_thumb_600_600_boston-city-flow.jpg"))

		err = os.Remove("/tmp/1_thumb_600_600_boston-city-flow.jpg")
		st.Expect(t, err, nil)
	}
}

//
// Test create File 02 - File too big.
//
func TestCreateFiles02(t *testing.T) {
	// test file.
	testFile := build.Default.GOPATH + "/src/app.skyclerk.com/backend/library/test/files/mr_93_e.pdf"

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Build file to post
	buffer, writer := buildLedgerFileform(t, testFile)

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

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"We have a 50MB upload limit."}`)
}

//
// Test create File 03 - Mime not supported.
//
func TestCreateFiles03(t *testing.T) {
	// test file.
	testFile := build.Default.GOPATH + "/src/app.skyclerk.com/backend/library/test/files/test01.txt"

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Build file to post
	buffer, writer := buildLedgerFileform(t, testFile)

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

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"We only allow image and pdf files to be uploaded."}`)
}

//
// Test create File 04 - Small PDF file with cases and spaces
//
func TestCreateFiles04(t *testing.T) {
	// test file.
	testFile := build.Default.GOPATH + "/src/app.skyclerk.com/backend/library/test/files/Income Statement copy.pdf"

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Build file to post
	buffer, writer := buildLedgerFileform(t, testFile)

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
	st.Expect(t, result.Name, "income-statement-copy.pdf")
	st.Expect(t, result.Type, "application/pdf")
	st.Expect(t, result.Size, int64(72689))
	st.Expect(t, true, strings.Contains(result.Url, "https://cdn-dev.skyclerk.com/accounts/33/1_income-statement-copy.pdf?Expires="))
	st.Expect(t, true, strings.Contains(result.Thumb600By600Url, "https://cdn-dev.skyclerk.com/accounts/33/1_thumb_600_600_income-statement-copy.jpeg?Expires="))

	// If we are testing locally (not on CI) we test to see if the file is on AWS with our signed key
	if len(os.Getenv("AWS_CLOUDFRONT_PRIVATE_SIGN_KEY")) > 0 {
		// Make sure the file is not already there.
		os.Remove("/tmp/1_income-statement-copy.pdf")

		err := downloadFile("/tmp/1_income-statement-copy.pdf", result.Url)
		st.Expect(t, err, nil)

		// Get the MD5 hash from the DB and compare
		ff := models.File{}
		db.New().Find(&ff, result.Id)
		st.Expect(t, ff.Hash, files.Md5("/tmp/1_income-statement-copy.pdf"))

		err = os.Remove("/tmp/1_income-statement-copy.pdf")
		st.Expect(t, err, nil)

		// --- here we chnage the key. We are just verifying AWS is honoring the singing

		err = downloadFile("/tmp/1_income-statement-copy.pdf", result.Url+"blah")
		st.Expect(t, err, nil)

		// Hash should be of "unauthorized"
		st.Expect(t, "b741482e554b40fd711a012fa74461cd", files.Md5("/tmp/1_income-statement-copy.pdf"))

		err = os.Remove("/tmp/1_income-statement-copy.pdf")
		st.Expect(t, err, nil)

		// ---- Test the thumb nail

		// Make sure the file is not already there.
		os.Remove("/tmp/1_thumb_600_600_income-statement-copy.jpeg")

		err = downloadFile("/tmp/1_thumb_600_600_income-statement-copy.jpeg", result.Thumb600By600Url)
		st.Expect(t, err, nil)
		st.Expect(t, "935cc90c4185d671c5f97fdb7d516baf", files.Md5("/tmp/1_thumb_600_600_income-statement-copy.jpeg"))

		err = os.Remove("/tmp/1_thumb_600_600_income-statement-copy.jpeg")
		st.Expect(t, err, nil)
	}
}

//
// Test create File 05 - Small PDF file
//
func TestCreateFiles05(t *testing.T) {
	// test file.
	testFile := build.Default.GOPATH + "/src/app.skyclerk.com/backend/library/test/files/apple.pdf"

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Build file to post
	buffer, writer := buildLedgerFileform(t, testFile)

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
	st.Expect(t, result.Name, "apple.pdf")
	st.Expect(t, result.Type, "application/pdf")
	st.Expect(t, result.Size, int64(95512))
	st.Expect(t, true, strings.Contains(result.Url, "https://cdn-dev.skyclerk.com/accounts/33/1_apple.pdf?Expires="))
}

// //
// // Test create File 06 - 35meg file. NOTE: this is commented out as it takes 30 seconds to run.
// //
// func TestCreateFiles06(t *testing.T) {
// 	// test file.
// 	testFile := build.Default.GOPATH + "/src/app.skyclerk.com/backend/library/test/files/Smiling-cowboy-standing-and-holding-lasso-519719714_7360x4912.jpeg"
//
// 	// Start the db connection.
// 	db, dbName, _ := models.NewTestDB("testing_db")
// 	defer models.TestingTearDown(db, dbName)
//
// 	// Create controller
// 	c := &Controller{}
// 	c.SetDB(db)
//
// 	// Build file to post
// 	buffer, writer := buildLedgerFileform(t, testFile)
//
// 	// Attach a ledger  to add this file to.
// 	err := writer.WriteField("ledger_id", "55")
// 	st.Expect(t, err, nil)
//
// 	// Setup request
// 	req, _ := http.NewRequest("POST", "/api/v3/33/files", buffer)
//
// 	// Set content type header
// 	req.Header.Set("Content-Type", writer.FormDataContentType())
//
// 	// Setup writer.
// 	w := httptest.NewRecorder()
// 	gin.SetMode("release")
// 	gin.DisableConsoleColor()
//
// 	r := gin.New()
// 	r.Use(func(c *gin.Context) {
// 		c.Set("accountId", 33)
// 		c.Set("userId", 109)
// 	})
// 	r.POST("/api/v3/:account/files", c.CreateFile)
// 	r.ServeHTTP(w, req)
//
// 	// Grab result and convert to strut
// 	result := models.File{}
// 	err = json.Unmarshal([]byte(w.Body.String()), &result)
//
// 	// Test results
// 	st.Expect(t, err, nil)
// 	st.Expect(t, w.Code, 201)
// 	st.Expect(t, result.Id, uint(1))
// 	st.Expect(t, result.AccountId, uint(33))
// 	st.Expect(t, result.Name, "smiling-cowboy-standing-and-holding-lasso-519719714-7360x4912.jpeg")
// 	st.Expect(t, result.Type, "image/jpeg")
// 	st.Expect(t, result.Size, int64(38737343))
// 	st.Expect(t, true, strings.Contains(result.Url, "https://cdn-dev.skyclerk.com/accounts/33/1_smiling-cowboy-standing-and-holding-lasso-519719714-7360x4912.jpeg?Expires="))
// 	st.Expect(t, true, strings.Contains(result.Thumb600By600Url, "https://cdn-dev.skyclerk.com/accounts/33/1_thumb_600_600_smiling-cowboy-standing-and-holding-lasso-519719714-7360x4912.jpeg?Expires="))
// }

//
// buildLedgerFileform so we can pust a file.
//
func buildLedgerFileform(t *testing.T, filePath string) (*bytes.Buffer, *multipart.Writer) {
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

	// Close writer
	err = writer.Close()
	st.Expect(t, err, nil)

	return body, writer
}

//
// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
//
func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

/* End File */
