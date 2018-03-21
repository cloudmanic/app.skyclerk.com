//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

type LabelsToLedger struct {
	LabelsToLedgerId        uint      `gorm:"primary_key;column:LabelsToLedgerId" json:"_"`
	LabelsToLedgerAccountId uint      `gorm:"column:LabelsToLedgerAccountId" json:"_"`
	LabelsToLedgerLabelId   uint      `gorm:"column:LabelsToLedgerLabelId" json:"_"`
	LabelsToLedgerLedgerId  uint      `gorm:"column:LabelsToLedgerLedgerId" json:"_"`
	LabelsToLedgerUpdatedAt time.Time `gorm:"column:LabelsToLedgerUpdatedAt" sql:"not null" json:"_"`
	LabelsToLedgerCreatedAt time.Time `gorm:"column:LabelsToLedgerCreatedAt" sql:"not null" json:"_"`
}

//
// Set the table name.
//
func (LabelsToLedger) TableName() string {
	return "LabelsToLedger"
}

/* End File */
