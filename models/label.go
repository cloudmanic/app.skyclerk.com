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

	// Make sure this category is not already in use.
	if action == "create" {

		c := Label{}

		if !db.New().Where("LabelsAccountId = ? AND LabelsName = ?", accountId, lbName).First(&c).RecordNotFound() {
			return errors.New(errMsg)
		}

		// Double check casing
		if strings.ToLower(lbName) == strings.ToLower(c.Name) {
			return errors.New(errMsg)
		}

	} else if action == "update" {

		c := Label{}

		if !db.New().Where("LabelsAccountId = ? AND LabelsName = ?", accountId, lbName).First(&c).RecordNotFound() {

			// Make sure it is not the same id as the one we are updating
			if c.Id != objId {
				return errors.New(errMsg)
			}
		}

		// Double check casing
		if (c.Id != objId) && (strings.ToLower(lbName) == strings.ToLower(c.Name)) {
			return errors.New(errMsg)
		}

	}

	// All good in the hood
	return nil
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
// DeleteLabelByAccountAndId - Delete a label by account and id.
//
func (db *DB) DeleteLabelByAccountAndId(accountId uint, labelId uint) error {
	// Make query to delete
	db.New().Where("LabelsAccountId = ? AND LabelsId = ?", accountId, labelId).Delete(Label{})

	// Delete from look up table.
	db.New().Where("LabelsToLedgerAccountId = ? AND LabelsToLedgerLabelId = ?", accountId, labelId).Delete(LabelsToLedger{})

	// Return result
	return nil
}

/* End File */
