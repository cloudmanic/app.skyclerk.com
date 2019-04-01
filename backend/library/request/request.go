//
// Date: 2018-03-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const defaultMysqlLimit = 100

//
// Get / Set standard query parms
//
func GetSetPagingParms(c *gin.Context) (int, int, int) {
	// Convert page to int.
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	// We do not allow limits over defaultMysqlLimit
	if limit > defaultMysqlLimit {
		limit = defaultMysqlLimit
	}

	if limit == 0 {
		limit = defaultMysqlLimit
	}

	// Offset can't be less than 0
	if offset < 0 {
		offset = 0
	}

	// Page can't be less than 1
	if page < 1 {
		page = 1
	}

	// Return happy.
	return page, limit, offset
}

/* End File */
