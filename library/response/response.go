//
// Date: 2018-03-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package response

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/gin-gonic/gin"
)

//
// Add paging info to the response.
//
func AddPagingInfoToHeaders(c *gin.Context, meta models.QueryMetaData) {
	c.Writer.Header().Set("X-Last-Page", strconv.FormatBool(meta.LastPage))
	c.Writer.Header().Set("X-Offset", strconv.Itoa(meta.Offset))
	c.Writer.Header().Set("X-Limit", strconv.Itoa(meta.Limit))
	c.Writer.Header().Set("X-No-Limit-Count", strconv.Itoa(meta.NoLimitCount))
}

//
// Results with meta in the header
//
func ResultsMeta(c *gin.Context, results interface{}, err error, meta models.QueryMetaData) {

	// Put meta data in header.
	AddPagingInfoToHeaders(c, meta)

	// Results
	Results(c, results, err)
}

//
// This is used when we are returning a list of results. Should
// almost never error. If no results are found it will be an empty array.
//
func Results(c *gin.Context, results interface{}, err error) {

	// Return json based on if this was a good result or not.
	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "There was an error. Please contact help@skyclerk.com for help."})

	} else {

		if c.DefaultQuery("output", "json") == "pretty" {
			c.String(http.StatusOK, interfaceToPretty(results))
		} else {
			c.JSON(http.StatusOK, results)
		}

	}

}

//
// Take an interface and return a pretty print version in json.
//
func interfaceToPretty(results interface{}) string {

	json, _ := json.Marshal(results)

	return jsonPrettyPrint(string(json))
}

//
// Take Json and and pretty print it.
//
func jsonPrettyPrint(in string) string {
	var out bytes.Buffer

	err := json.Indent(&out, []byte(in), "", "\t")

	if err != nil {
		return in
	}

	return out.String()
}

/* End File */
