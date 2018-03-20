//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: spicer
// Last Modified: 2018-03-20
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

type Label struct {
	Id        uint      `gorm:"primary_key;column:LabelsId" json:"id"`
	AccountId uint      `gorm:"column:LabelsAccountId" sql:"not null" json:"account_id"`
	UpdatedAt time.Time `gorm:"column:LabelsUpdatedAt" sql:"not null" json:"_"`
	CreatedAt time.Time `gorm:"column:LabelsCreatedAt" sql:"not null" json:"_"`
	Name      string    `gorm:"column:LabelsName" sql:"not null;" json:"name"`
	System    uint      `gorm:"column:LabelsSystem" sql:"not null" json:"_"`
}

//
// Set the table name.
//
func (Label) TableName() string {
	return "Labels"
}

/* End File */
