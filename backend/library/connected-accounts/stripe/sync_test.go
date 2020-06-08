//
// Date: 2020-06-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

import (
	"testing"

	"app.skyclerk.com/backend/models"
)

//
// TestSync01 will sync stripe
//
func TestSync01(t *testing.T) {
	// Connected account model.
	ac := models.ConnectedAccounts{
		AccountID:            33,
		Connection:           "Stripe",
		StripeUserID:         "Ris6D3eYxtXbIkRz5K7aJBLiGkTOvuuD",
		StripeAccessToken:    "sk_test_1Ris6D3eYxtXbIkRz5K7aJBLiGkTOvuuDDw9UFikFsw47q8zNruW0fCDpZfuqwdriHNnKiJFFKEM6yNEDeqKzKNgK007Gn0O8HV",
		StripeRefreshToken:   "rt_HQUXQqXe9KvzfV60vfFnNqVnGNSwBMtpXGUr4N0YrCc2Feqq",
		StripeScope:          "",
		StripeLastItem:       1487731749,
		StripePublishableKey: "pk_test_1Ris6D3eYxtXbIkRz5K7aJBLiGkTOvuuDzkgPVNu7mF1aON92g9n6xREVSPcF134HGoKuuh9sCwgt8Ai1D6ApFaR100JiHDosXn",
	}

	Sync(ac)

}

/* End File */
