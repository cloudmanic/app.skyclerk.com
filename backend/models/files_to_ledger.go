//
// Date: 2019-04-21
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

type FilesToLedger struct {
	FilesToLedgerFileId   uint `gorm:"column:FilesToLedgerFileId" json:"_"`
	FilesToLedgerLedgerId uint `gorm:"column:FilesToLedgerLedgerId" json:"_"`
}

//
// Set the table name.
//
func (FilesToLedger) TableName() string {
	return "FilesToLedger"
}

/* End File */
