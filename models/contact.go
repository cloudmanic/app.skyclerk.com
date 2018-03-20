//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: spicer
// Last Modified: 2018-03-20
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

type Contact struct {
	Id            uint      `gorm:"primary_key;column:ContactsId" json:"id"`
	AccountId     uint      `gorm:"column:ContactsAccountId" sql:"not null" json:"account_id"`
	UpdatedAt     time.Time `gorm:"column:ContactsUpdatedAt" sql:"not null" json:"created_at"`
	CreatedAt     time.Time `gorm:"column:ContactsCreatedAt" sql:"not null" json:"updated_at"`
	Name          string    `gorm:"column:ContactsName" sql:"not null" json:"name"`
	FirstName     string    `gorm:"column:ContactsFirstName" sql:"not null" json:"first"`
	LastName      string    `gorm:"column:ContactsLastName" sql:"not null" json:"name"`
	Address       string    `gorm:"column:ContactsAddress" sql:"not null" json:"_"`
	City          string    `gorm:"column:ContactsCity" sql:"not null" json:"_"`
	State         string    `gorm:"column:ContactsState" sql:"not null" json:"_"`
	Zip           string    `gorm:"column:ContactsZip" sql:"not null" json:"_"`
	Phone         string    `gorm:"column:ContactsPhone" sql:"not null" json:"_"`
	Fax           string    `gorm:"column:ContactsFax" sql:"not null" json:"_"`
	Website       string    `gorm:"column:ContactsWebsite" sql:"not null" json:"_"`
	AccountNumber string    `gorm:"column:ContactsAccountNumber" sql:"not null" json:"_"`
	Avatar        string    `gorm:"column:ContactsAvatar" sql:"not null" json:"_"`
	AvatarChecked string    `gorm:"column:ContactsAvatarChecked" sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"_"`
	Email         string    `gorm:"column:ContactsEmail" sql:"not null" json:"_"`
	Twitter       string    `gorm:"column:ContactsTwitter" sql:"not null" json:"_"`
	Facebook      string    `gorm:"column:ContactsFacebook" sql:"not null" json:"_"`
	Linkedin      string    `gorm:"column:ContactsLinkedin" sql:"not null" json:"_"`
	Type          string    `gorm:"column:ContactsType" sql:"not null;type:ENUM('Customer', 'Vendor', 'Both');default:'Both'" json:"_"`
	HrId          uint64    `gorm:"column:ContactsHrId" sql:"not null" json:"_"`
	PricingPlanId uint      `gorm:"column:ContactsPricingPlanId" sql:"not null" json:"_"`
	GatewayId     uint      `gorm:"column:ContactsGatewayId" sql:"not null" json:"_"`
	CardMask      string    `gorm:"column:ContactsCardMask" sql:"not null" json:"_"`
	CardType      string    `gorm:"column:ContactsCardType" sql:"not null" json:"_"`
	CardExpire    string    `gorm:"column:ContactsCardExpire" sql:"not null" json:"_"`
	Country       string    `gorm:"column:ContactsCountry" sql:"not null" json:"_"`
}

//
// Set the table name.
//
func (Contact) TableName() string {
	return "Contacts"
}

/* End File */
