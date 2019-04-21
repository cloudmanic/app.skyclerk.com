//
// Date: 2019-04-21
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

type FilesToLedger struct {
	FilesToLedgerId        uint      `gorm:"primary_key;column:FilesToLedgerId" json:"_"`
	FilesToLedgerAccountId uint      `gorm:"column:FilesToLedgerAccountId" json:"_"`
	FilesToLedgerFileId    uint      `gorm:"column:FilesToLedgerFileId" json:"_"`
	FilesToLedgerLedgerId  uint      `gorm:"column:FilesToLedgerLedgerId" json:"_"`
	FilesToLedgerUpdatedAt time.Time `gorm:"column:FilesToLedgerUpdatedAt" sql:"not null" json:"_"`
	FilesToLedgerCreatedAt time.Time `gorm:"column:FilesToLedgerCreatedAt" sql:"not null" json:"_"`
}

//
// Set the table name.
//
func (FilesToLedger) TableName() string {
	return "FilesToLedger"
}

/* End File */
