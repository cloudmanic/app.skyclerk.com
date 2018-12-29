//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-29
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"
)

type Label struct {
	Id        uint      `gorm:"primary_key;column:LabelsId" json:"id"`
	AccountId uint      `gorm:"column:LabelsAccountId" sql:"not null" json:"_"`
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

//
// Return a label by account and id.
//
func (db *DB) GetLabelByAccountAndId(accountId uint, labelId uint) (Label, error) {

	l := Label{}

	// Make query
	if db.New().Where("LabelsAccountId = ? AND LabelsId = ?", accountId, labelId).First(&l).RecordNotFound() {
		return Label{}, errors.New("Label not found.")
	}

	// Return result
	return l, nil
}

/* End File */
