//
// Date: 2019-03-14
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"testing"

	"github.com/cloudmanic/skyclerk.com/library/test"
	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/davecgh/go-spew/spew"
)

//
// TestGetLedgers01 Test get ledgers 01
//
func TestGetLedgers01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	l1 := test.GetRandomLedger()
	db.LedgerCreate(&l1)

	spew.Dump(l1)

	// // Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	// db.Save(&models.Contact{AccountId: 34, Name: "Apple Inc.", FirstName: "", LastName: ""})
	// db.Save(&models.Contact{AccountId: 34, Name: "Matthews Etc. LLC", FirstName: "Spicer", LastName: "Matthews"})
	// db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})
	// db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Mike", LastName: "Rosso"})
	// db.Save(&models.Contact{AccountId: 33, Name: "Zoo Inc.", FirstName: "", LastName: ""})
	// db.Save(&models.Contact{AccountId: 33, Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"})
	// db.Save(&models.Contact{AccountId: 33, Name: "Dope Dealer, LLC", FirstName: "", LastName: ""})
	//
	// // Setup request
	// req, _ := http.NewRequest("GET", "/api/v1/33/contacts", nil)
	//
	// // Setup writer.
	// w := httptest.NewRecorder()
	// gin.SetMode("release")
	// gin.DisableConsoleColor()
	//
	// r := gin.New()
	// r.Use(func(c *gin.Context) {
	// 	c.Set("accountId", 33)
	// 	c.Set("userId", uint(109))
	// })
	// r.GET("/api/v1/:account/contacts", c.GetContacts)
	// r.ServeHTTP(w, req)
	//
	// // Grab result and convert to strut
	// results := []models.Contact{}
	// err := json.Unmarshal([]byte(w.Body.String()), &results)
	//
	// // Test results
	// st.Expect(t, err, nil)
	// st.Expect(t, results[0].Id, uint(3))
	// st.Expect(t, results[1].Id, uint(4))
	// st.Expect(t, results[2].Id, uint(6))
	// st.Expect(t, results[3].Id, uint(7))
	// st.Expect(t, results[4].Id, uint(5))
	// st.Expect(t, results[0].Name, "")
	// st.Expect(t, results[1].Name, "")
	// st.Expect(t, results[2].Name, "Abc Inc.")
	// st.Expect(t, results[3].Name, "Dope Dealer, LLC")
	// st.Expect(t, results[4].Name, "Zoo Inc.")
}

/* End File */
