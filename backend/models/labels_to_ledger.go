//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

type LabelsToLedger struct {
	LabelsToLedgerLabelId  uint `gorm:"column:LabelsToLedgerLabelId" json:"_"`
	LabelsToLedgerLedgerId uint `gorm:"column:LabelsToLedgerLedgerId" json:"_"`
}

//
// Set the table name.
//
func (LabelsToLedger) TableName() string {
	return "LabelsToLedger"
}

/* End File */
