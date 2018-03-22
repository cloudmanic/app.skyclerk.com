//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

type Category struct {
	Id        uint      `gorm:"primary_key;column:CategoriesId" json:"id"`
	AccountId uint      `gorm:"column:CategoriesAccountId" sql:"not null" json:"_"`
	UpdatedAt time.Time `gorm:"column:CategoriesUpdatedAt" sql:"not null" json:"_"`
	CreatedAt time.Time `gorm:"column:CategoriesCreatedAt" sql:"not null" json:"_"`
	Name      string    `gorm:"column:CategoriesName" sql:"not null;" json:"name"`
	Type      string    `gorm:"column:CategoriesType" sql:"not null" json:"type"`
	Irs       string    `gorm:"column:CategoriesIrs" sql:"not null" json:"_"`
	Show      string    `gorm:"column:CategoriesShow" sql:"not null" json:"_"`
}

//
// Set the table name.
//
func (Category) TableName() string {
	return "Categories"
}

/* End File */
