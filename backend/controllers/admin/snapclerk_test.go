//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"bytes"
	"encoding/json"
	"go/build"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestConvertSnapClerk01 - test converting a snapclerk to a ledger entry.
//
func TestConvertSnapClerk01(t *testing.T) {
	// test file.
	destinationFile := "/tmp/money-2724241_1920.jpg"
	sourceFile := build.Default.GOPATH + "/src/app.skyclerk.com/backend/library/test/files/money-2724241_1920.jpg"

	// Make a copy of test file
	input, err := ioutil.ReadFile(sourceFile)
	st.Expect(t, err, nil)

	err = ioutil.WriteFile(destinationFile, input, 0644)
	st.Expect(t, err, nil)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)
	db.Save(&models.AcctToUsers{AccountId: uint(33), UserId: user.Id})

	account := test.GetRandomAccount(33)
	account.OwnerId = user.Id
	account.Name = "Matthews Etc."
	db.Save(&account)
	user.Accounts = append(user.Accounts, account)

	// Store a file
	file, err := db.StoreFile(account.Id, destinationFile)
	st.Expect(t, err, nil)

	// Post data
	lPost := models.SnapClerk{
		FileId:    file.Id,
		AccountId: 33,
		AddedById: 1,
		Labels:    "hahahah , thanks",
		Note:      "4 inch drain for inspection on 11/7",
		Lat:       "45.28819058891675",
		Lon:       "-122.93470961648806",
	}

	// Save the snapclerk in db.
	db.New().Save(&lPost)

	// This is what we would add on in the FE.
	lPost.Contact = "Home Depot"
	lPost.Amount = 88.22
	lPost.Category = "Travel"

	// Get JSON
	postStr, _ := json.Marshal(lPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/snapclerk/convert/1", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/snapclerk/convert/:id", c.ConvertSnapClerk)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Ledger{}
	err = json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Amount, (lPost.Amount * -1))
	st.Expect(t, result.Note, lPost.Note)
	st.Expect(t, result.Lat, 45.28819058891675)
	st.Expect(t, result.Lon, -122.93470961648806)
	st.Expect(t, result.Contact.Name, lPost.Contact)
	st.Expect(t, result.Category.Name, lPost.Category)
	st.Expect(t, result.Labels[0].Name, "Snap!Clerk")
	st.Expect(t, result.Labels[1].Name, "hahahah")
	st.Expect(t, result.Labels[2].Name, "thanks")
	st.Expect(t, true, strings.Contains(result.Files[0].Url, "https://cdn-dev.skyclerk.com/accounts/33/1_money-2724241-1920.jpg?Expires="))

	// Get ledger.
	l, err := db.GetLedgerByAccountAndId(uint(33), result.Id)
	st.Expect(t, err, nil)

	// Double check the db.
	st.Expect(t, l.Amount, (lPost.Amount * -1))
	st.Expect(t, l.Note, lPost.Note)
	st.Expect(t, l.Lat, 45.28819058891675)
	st.Expect(t, l.Lon, -122.93470961648806)
	st.Expect(t, l.Contact.Name, lPost.Contact)
	st.Expect(t, l.Category.Name, lPost.Category)
	st.Expect(t, l.Labels[0].Name, "Snap!Clerk")
	st.Expect(t, l.Labels[1].Name, "hahahah")
	st.Expect(t, l.Labels[2].Name, "thanks")
	st.Expect(t, true, strings.Contains(l.Files[0].Url, "https://cdn-dev.skyclerk.com/accounts/33/1_money-2724241-1920.jpg?Expires="))

	// Double check activity
	ac := models.Activity{}
	db.New().Find(&ac, 1)
	st.Expect(t, ac.Action, "expense")
	st.Expect(t, ac.Name, "Home Depot")
	st.Expect(t, ac.LedgerId, l.Id)

	// Double check snapclerk
	s := models.SnapClerk{}
	db.New().Find(&s, 1)
	st.Expect(t, s.ReviewedById, uint(109))
	st.Expect(t, s.LedgerId, uint(1))
	st.Expect(t, s.Status, "Processed")
	st.Expect(t, s.Amount, (lPost.Amount * -1))
}

/* End File */
