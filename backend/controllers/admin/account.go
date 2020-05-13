//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Account struct
type Account struct {
	AccountID    uint      `json:"account_id"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
	Name         string    `json:"name"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Status       string    `json:"status"`
	LedgerCount  int       `json:"ledger_count"`
}

//
// GetAccounts returns a list of accounts.
//
func (t *Controller) GetAccounts(c *gin.Context) {
	// Raw query
	sql := `
SELECT
	accounts.id AS account_id,
	accounts.created_at AS created_at,
	accounts.last_activity AS last_activity,
	accounts.name AS name,
	users.first_name AS first_name,
	users.last_name AS last_name,
	users.email AS email,
	billings.status AS status,
	count(LedgerId) AS ledger_count
FROM
	Ledger
	INNER JOIN accounts ON Ledger.LedgerAccountId = accounts.Id
	INNER JOIN users ON accounts.owner_id = users.id
	INNER JOIN billings ON accounts.billing_id = billings.id
GROUP BY
	accounts.Id
ORDER BY
	accounts.last_activity DESC,
	accounts.created_at DESC
LIMIT 100`

	rt := []Account{}

	// Run query.
	t.db.New().Raw(sql).Scan(&rt)

	// Return happy JSON
	c.JSON(200, rt)
}

/* End File */
