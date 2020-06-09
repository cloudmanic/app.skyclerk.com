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

// CategoryUsage struct
type CategoryUsage struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Category struct
type Category struct {
	Id        uint      `gorm:"primary_key;column:CategoriesId" json:"id"`
	AccountId uint      `gorm:"column:CategoriesAccountId" sql:"not null" json:"account_id"`
	UpdatedAt time.Time `gorm:"column:CategoriesUpdatedAt" sql:"not null" json:"-"`
	CreatedAt time.Time `gorm:"column:CategoriesCreatedAt" sql:"not null" json:"-"`
	Name      string    `gorm:"column:CategoriesName" sql:"not null;" json:"name"`
	Type      string    `gorm:"column:CategoriesType" sql:"not null" json:"type"` // 1 = expense, 2 = income
	Irs       string    `gorm:"column:CategoriesIrs" sql:"not null" json:"-"`
	Show      string    `gorm:"column:CategoriesShow" sql:"not null" json:"-"`
	Count     int       `gorm:"-" sql:"not null" json:"count"`
}

//
// Set the table name.
//
func (a Category) TableName() string {
	return "Categories"
}

//
// Validate for this model.
//
func (a Category) Validate(db Datastore, action string, userId uint, accountId uint, objId uint) error {
	return validation.ValidateStruct(&a,

		validation.Field(&a.Name,
			validation.Required.Error("The name field is required."),
			validation.By(func(value interface{}) error { return db.ValidateDuplicateCategoryName(a, accountId, objId, action) }),
		),

		validation.Field(&a.Type,
			validation.Required.Error("The type field is required."),
			validation.In("1", "2").Error("The type field must be 1, or 2."),
		),
	)
}

//
// ValidateDuplicateCategoryName - Validate Duplicate Name
//
func (db *DB) ValidateDuplicateCategoryName(cat Category, accountId uint, objId uint, action string) error {

	const errMsg = "Category name is already in use."

	// trim any white space
	catName := strings.Trim(cat.Name, " ")

	// Make sure this category is not already in use.
	if action == "create" {

		c := Category{}

		if !db.New().Where("CategoriesAccountId = ? AND CategoriesName = ? AND CategoriesType = ?", accountId, catName, cat.Type).First(&c).RecordNotFound() {
			return errors.New(errMsg)
		}

		// Double check casing
		if strings.ToLower(catName) == strings.ToLower(c.Name) {
			return errors.New(errMsg)
		}

	} else if action == "update" {

		c := Category{}

		if !db.New().Where("CategoriesAccountId = ? AND CategoriesName = ? AND CategoriesType = ?", accountId, catName, cat.Type).First(&c).RecordNotFound() {
			// Make sure it is not the same id as the one we are updating
			if c.Id != objId {
				return errors.New(errMsg)
			}
		}

		// Double check casing
		if (c.Id != objId) && (strings.ToLower(catName) == strings.ToLower(c.Name)) {
			return errors.New(errMsg)
		}

	}

	// All good in the hood
	return nil
}

//
// Return a category by account and id.
//
func (db *DB) GetCategoryByAccountAndId(accountId uint, categoryId uint) (Category, error) {

	c := Category{}

	// Make query
	if db.New().Where("CategoriesAccountId = ? AND CategoriesId = ?", accountId, categoryId).First(&c).RecordNotFound() {
		return Category{}, errors.New("Category not found.")
	}

	// Return result
	return c, nil
}

//
// GetCategoryByNameAndTypeAndAccountID returns a category by name and account id.
//
func (db *DB) GetCategoryByNameAndTypeAndAccountID(accountID uint, name string, catType string) (Category, error) {
	c := Category{}

	// Make query
	if db.New().Where("CategoriesAccountId = ? AND CategoriesName = ? AND CategoriesType = ?", accountID, name, catType).First(&c).RecordNotFound() {
		return Category{}, errors.New("Category not found.")
	}

	// Return result
	return c, nil
}

//
// GetOrCreateCategory will get or create category if we do not already have.
//
func (db *DB) GetOrCreateCategory(accountID uint, name string, catType string) Category {
	// See if we already have the category.
	cat, err := db.GetCategoryByNameAndTypeAndAccountID(accountID, name, catType)

	if err != nil {
		cat.Type = catType
		cat.Name = name
		cat.AccountId = accountID
		db.New().Save(&cat)
	}

	return cat
}

//
// DeleteCategoryByAccountAndId - Delete a category by account and id.
//
func (db *DB) DeleteCategoryByAccountAndId(accountId uint, categoryId uint) error {

	// Make query to see if we have ledger entries with this category
	if !db.New().Where("LedgerAccountId = ? AND LedgerCategoryId = ?", accountId, categoryId).First(&Ledger{}).RecordNotFound() {
		return errors.New("Can not delete category. It is in use by a ledger entry.")
	}

	// Make query
	db.New().Where("CategoriesAccountId = ? AND CategoriesId = ?", accountId, categoryId).Delete(Category{})

	// Return result
	return nil
}

//
// GetCategoryUsage - returns a list of categories by account and the usage.
//
func (db *DB) GetCategoryUsageByAccount(accountId uint) []CategoryUsage {
	// SQL String
	sql := "SELECT CategoriesName AS name, COUNT(LedgerId) AS count FROM Ledger "
	sql = sql + "INNER JOIN Categories ON Ledger.LedgerCategoryId=Categories.CategoriesId "
	sql = sql + "WHERE CategoriesAccountId = ? "
	sql = sql + "GROUP BY CategoriesName ORDER BY CategoriesName "

	// Struct we return
	rt := []CategoryUsage{}

	// Run query
	db.New().Raw(sql, accountId).Scan(&rt)

	// Return happy.
	return rt
}

//
// LoadDefaultCategories - install the default categories we get on a new account.
//
func (db *DB) LoadDefaultCategories(accountId uint) {
	// Default cats
	cats := []Category{
		// Expenses
		{Name: "Advertising", Type: "1", AccountId: accountId},
		{Name: "Car & Truck Expenses", Type: "1", AccountId: accountId},
		{Name: "Commissions & Fees", Type: "1", AccountId: accountId},
		{Name: "Insurance", Type: "1", AccountId: accountId},
		{Name: "Mortgage", Type: "1", AccountId: accountId},
		{Name: "Meals & Entertainment", Type: "1", AccountId: accountId},
		{Name: "Office Expense", Type: "1", AccountId: accountId},
		{Name: "Professional Services", Type: "1", AccountId: accountId},
		{Name: "Supplies", Type: "1", AccountId: accountId},
		{Name: "Travel", Type: "1", AccountId: accountId},
		{Name: "Maintenance", Type: "1", AccountId: accountId},
		{Name: "Contractors & Freelancers", Type: "1", AccountId: accountId},
		{Name: "Cost of Goods Sold'", Type: "1", AccountId: accountId},
		{Name: "Equipment Rental", Type: "1", AccountId: accountId},
		{Name: "Utilities", Type: "1", AccountId: accountId},
		{Name: "Employee Wage", Type: "1", AccountId: accountId},
		{Name: "Taxes and Licenses", Type: "1", AccountId: accountId},
		{Name: "Pension & Profit-Sharing Plans", Type: "1", AccountId: accountId},
		{Name: "Rent Or Lease'", Type: "1", AccountId: accountId},
		{Name: "Other Expense", Type: "1", AccountId: accountId},

		// income
		{Name: "Sales", Type: "2", AccountId: accountId},
		{Name: "Returns", Type: "2", AccountId: accountId},
		{Name: "Other Income", Type: "2", AccountId: accountId},
	}

	// Save to database
	for _, row := range cats {
		db.New().Create(&row)
	}
}

/* End File */
