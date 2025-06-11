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
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// LabelUsage struct
type LabelUsage struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Label struct
type Label struct {
	Id        uint      `gorm:"primary_key;column:LabelsId" json:"id"`
	AccountId uint      `gorm:"column:LabelsAccountId" sql:"not null" json:"account_id"`
	UpdatedAt time.Time `gorm:"column:LabelsUpdatedAt" sql:"not null" json:"_"`
	CreatedAt time.Time `gorm:"column:LabelsCreatedAt" sql:"not null" json:"_"`
	Name      string    `gorm:"column:LabelsName" sql:"not null;" json:"name"`
	System    uint      `gorm:"column:LabelsSystem" sql:"not null" json:"_"`
	Count     int       `gorm:"-" sql:"not null" json:"count"`
}

//
// Set the table name.
//
func (Label) TableName() string {
	return "Labels"
}

//
// Validate for this model.
//
func (a Label) Validate(db Datastore, action string, userId uint, accountId uint, objId uint) error {
	return validation.ValidateStruct(&a,

		validation.Field(&a.Name,
			validation.Required.Error("The name field is required."),
			validation.By(func(value interface{}) error { return db.ValidateDuplicateLabelName(a, accountId, objId, action) }),
		),
	)
}

//
// ValidateDuplicateLabelName - Validate Duplicate Name
//
func (db *DB) ValidateDuplicateLabelName(obj Label, accountId uint, objId uint, action string) error {

	const errMsg = "Label name is already in use."

	// trim any white space
	lbName := strings.Trim(obj.Name, " ")

	// Make sure this label is not already in use.
	if action == "create" {
		var labels []Label
		
		// Find labels with similar names (case-insensitive check)
		db.New().Where("LabelsAccountId = ?", accountId).Find(&labels)
		
		for _, c := range labels {
			if strings.ToLower(strings.Trim(c.Name, " ")) == strings.ToLower(lbName) {
				return errors.New(errMsg)
			}
		}

	} else if action == "update" {
		var labels []Label
		
		// Find labels with similar names (case-insensitive check)
		db.New().Where("LabelsAccountId = ? AND LabelsId != ?", accountId, objId).Find(&labels)
		
		for _, c := range labels {
			if strings.ToLower(strings.Trim(c.Name, " ")) == strings.ToLower(lbName) {
				return errors.New(errMsg)
			}
		}
	}

	// All good in the hood
	return nil
}

//
// GetOrCreateLabel get or create a label
//
func (db *DB) GetOrCreateLabel(accountID uint, name string) Label {
	// Get label
	label, err := db.GetLabelByAccountAndName(accountID, name)

	// No Label found let's create.
	if err != nil {
		label.Name = name
		label.AccountId = accountID
		db.New().Save(&label)
	}

	return label
}

//
// GetLabelByAccountAndId - Return a label by account and id.
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

//
// GetLabelByAccountAndName - Return a label by account and name.
//
func (db *DB) GetLabelByAccountAndName(accountId uint, name string) (Label, error) {
	l := Label{}

	// Make query
	if db.New().Where("LabelsAccountId = ? AND LabelsName = ?", accountId, name).First(&l).RecordNotFound() {
		return Label{}, errors.New("Label not found.")
	}

	// Return result
	return l, nil
}

//
// DeleteLabelByAccountAndId - Delete a label by account and id.
//
func (db *DB) DeleteLabelByAccountAndId(accountId uint, labelId uint) error {
	// Make query to delete
	db.New().Where("LabelsAccountId = ? AND LabelsId = ?", accountId, labelId).Delete(Label{})

	// Delete from look up table.
	db.New().Where("LabelsToLedgerLabelId = ?", labelId).Delete(LabelsToLedger{})

	// Return result
	return nil
}

//
// GetLabelUsageByAccount - returns a list of labels by account and the usage.
//
func (db *DB) GetLabelUsageByAccount(accountId uint) []LabelUsage {
	// SQL String
	sql := "SELECT LabelsName as name, COUNT(LabelsToLedgerLedgerId) as count FROM LabelsToLedger "
	sql = sql + "INNER JOIN Labels ON LabelsToLedger.LabelsToLedgerLabelId=Labels.LabelsId "
	sql = sql + "WHERE LabelsAccountId = ? "
	sql = sql + "GROUP BY LabelsName ORDER BY LabelsName"

	// Struct we return
	rt := []LabelUsage{}

	// Run query
	db.New().Raw(sql, accountId).Scan(&rt)

	// Return happy.
	return rt
}

/* End File */
