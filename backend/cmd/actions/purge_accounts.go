//
// Date: 2025-01-11
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2025-01-11
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package actions

import (
	"fmt"
	"time"

	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

type PurgeStats struct {
	TotalAccountsChecked    int
	AccountsWithNoLedgers   int
	AccountsWithOldLedgers  int
	AccountsDeleted         int
	AccountsSkippedForEmail int
	Errors                  []string
}

//
// PurgeOldAccounts - Delete accounts that meet certain criteria:
// 1. Accounts with zero ledger entries created over 6 months ago
// 2. Accounts that haven't had a ledger entry in one year
// 3. Skip any account associated with spicer@cloudmanic.com
//
func PurgeOldAccounts(db models.Datastore) {
	fmt.Println("Starting account purge process...")
	fmt.Println("================================================")

	stats := PurgeStats{
		Errors: []string{},
	}

	// Get all accounts
	accounts := []models.Account{}
	db.New().Find(&accounts)
	stats.TotalAccountsChecked = len(accounts)

	// Time boundaries
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	for _, account := range accounts {
		// Check if account is associated with spicer@cloudmanic.com
		users := db.GetUsersByAccount(account.Id)
		skipAccount := false
		for _, user := range users {
			if user.Email == "spicer@cloudmanic.com" {
				skipAccount = true
				stats.AccountsSkippedForEmail++
				break
			}
		}

		if skipAccount {
			continue
		}

		// Get ledger entries for this account
		ledgers := []models.Ledger{}
		db.New().Where("LedgerAccountId = ?", account.Id).Order("LedgerCreatedAt DESC").Find(&ledgers)

		shouldDelete := false
		deleteReason := ""

		if len(ledgers) == 0 {
			// Account has zero ledger entries
			if account.CreatedAt.Before(sixMonthsAgo) {
				shouldDelete = true
				deleteReason = "zero ledgers and created over 6 months ago"
				stats.AccountsWithNoLedgers++
			}
		} else {
			// Check if the most recent ledger entry is older than one year
			mostRecentLedger := ledgers[0]
			if mostRecentLedger.CreatedAt.Before(oneYearAgo) {
				shouldDelete = true
				deleteReason = "no ledger activity in over one year"
				stats.AccountsWithOldLedgers++
			}
		}

		if shouldDelete {
			fmt.Printf("Deleting account %d (%s) - Reason: %s\n", account.Id, account.Name, deleteReason)
			
			// Clear the account data first
			db.ClearAccount(account.Id)
			
			// Delete the account and associated records
			db.DeleteAccount(account.Id)
			
			stats.AccountsDeleted++
			
			// Log the deletion
			services.InfoMsg(fmt.Sprintf("Purged account %d (%s) - Reason: %s", account.Id, account.Name, deleteReason))
		}
	}

	// Print summary statistics
	fmt.Println("\n================================================")
	fmt.Println("Account Purge Summary")
	fmt.Println("================================================")
	fmt.Printf("Total accounts checked:        %d\n", stats.TotalAccountsChecked)
	fmt.Printf("Accounts with zero ledgers:    %d\n", stats.AccountsWithNoLedgers)
	fmt.Printf("Accounts with old ledgers:     %d\n", stats.AccountsWithOldLedgers)
	fmt.Printf("Accounts deleted:              %d\n", stats.AccountsDeleted)
	fmt.Printf("Accounts skipped (email rule): %d\n", stats.AccountsSkippedForEmail)
	fmt.Printf("Errors encountered:            %d\n", len(stats.Errors))
	
	if len(stats.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range stats.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}
	
	fmt.Println("\nPurge process completed successfully!")
}

/* End File */