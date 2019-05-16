//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2019-01-13
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"app.skyclerk.com/backend/library/avatar"
	"app.skyclerk.com/backend/library/files"
	"app.skyclerk.com/backend/library/store/object"
	"app.skyclerk.com/backend/services"

	"github.com/eefret/gravatar"
	validation "github.com/go-ozzo/ozzo-validation"
)

// Conact struct
type Contact struct {
	Id            uint      `gorm:"primary_key;column:ContactsId" json:"id"`
	AccountId     uint      `gorm:"column:ContactsAccountId" sql:"not null" json:"account_id"`
	UpdatedAt     time.Time `gorm:"column:ContactsUpdatedAt" sql:"not null" json:"_"`
	CreatedAt     time.Time `gorm:"column:ContactsCreatedAt" sql:"not null" json:"_"`
	Name          string    `gorm:"column:ContactsName" sql:"not null" json:"name"`
	FirstName     string    `gorm:"column:ContactsFirstName" sql:"not null" json:"first_name"`
	LastName      string    `gorm:"column:ContactsLastName" sql:"not null" json:"last_name"`
	Address       string    `gorm:"column:ContactsAddress" sql:"not null" json:"address"`
	City          string    `gorm:"column:ContactsCity" sql:"not null" json:"city"`
	State         string    `gorm:"column:ContactsState" sql:"not null" json:"state"`
	Zip           string    `gorm:"column:ContactsZip" sql:"not null" json:"zip"`
	Phone         string    `gorm:"column:ContactsPhone" sql:"not null" json:"phone"`
	Fax           string    `gorm:"column:ContactsFax" sql:"not null" json:"fax"`
	Website       string    `gorm:"column:ContactsWebsite" sql:"not null" json:"website"`
	AccountNumber string    `gorm:"column:ContactsAccountNumber" sql:"not null" json:"account_number"`
	Avatar        string    `gorm:"column:ContactsAvatar" sql:"not null" json:"_"`
	AvatarChecked string    `gorm:"column:ContactsAvatarChecked" sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"_"` // This means someone has not uploaded an image we have just generated on. So we can update it if we want.
	AvatarUrl     string    `gorm:"-" json:"avatar_url"`                                                                      // Not stored in DB.
	Email         string    `gorm:"column:ContactsEmail" sql:"not null" json:"email"`
	Twitter       string    `gorm:"column:ContactsTwitter" sql:"not null" json:"twitter"`
	Facebook      string    `gorm:"column:ContactsFacebook" sql:"not null" json:"facebook"`
	Linkedin      string    `gorm:"column:ContactsLinkedin" sql:"not null" json:"linkedin"`
	Type          string    `gorm:"column:ContactsType" sql:"not null;type:ENUM('Customer', 'Vendor', 'Both');default:'Both'" json:"_"`
	HrId          uint64    `gorm:"column:ContactsHrId" sql:"not null" json:"_"`
	PricingPlanId uint      `gorm:"column:ContactsPricingPlanId" sql:"not null" json:"_"`
	GatewayId     uint      `gorm:"column:ContactsGatewayId" sql:"not null" json:"_"`
	CardMask      string    `gorm:"column:ContactsCardMask" sql:"not null" json:"_"`
	CardType      string    `gorm:"column:ContactsCardType" sql:"not null" json:"_"`
	CardExpire    string    `gorm:"column:ContactsCardExpire" sql:"not null" json:"_"`
	Country       string    `gorm:"column:ContactsCountry" sql:"not null" json:"country"`
}

// generateAvatarsWorkerJob struct
type generateAvatarsWorkerJob struct {
	name      string
	count     int
	rowCount  int
	accountId uint
	id        uint
	email     string
}

//
// Set the table name.
//
func (Contact) TableName() string {
	return "Contacts"
}

//
// Validate for this model.
//
func (a Contact) Validate(db Datastore, action string, userId uint, accountId uint, objId uint) error {

	return validation.ValidateStruct(&a,

		validation.Field(&a.Name,
			validation.By(func(value interface{}) error { return db.ValidateContactNameOrFirstLast(a, accountId, objId, action) }),
		),
	)
}

//
// CreateContact
//
func (db *DB) CreateContact(contact *Contact) error {
	// Run query to save
	db.Create(&contact)

	// Figure out name
	name := contact.Name

	if len(contact.Name) == 0 {
		name = fmt.Sprintf("%s %s", contact.FirstName, contact.LastName)
	}

	// Generate and Store avatar
	avatarPath, err := GenerateAndStoreAvatar(contact.AccountId, contact.Id, name, contact.Email)

	if err != nil {
		return err
	}

	// Update DB with the new path of avatar
	contact.Avatar = avatarPath
	db.Save(&contact)

	// Add a signed avatar path
	contact.AvatarUrl = db.GetSignedFileUrl(contact.Avatar)

	// Return happy
	return nil
}

//
// Validate Duplicate Contact Name, First, Last
//
func (db *DB) ValidateContactNameOrFirstLast(contact Contact, accountId uint, objId uint, action string) error {

	// trim any white space
	contactName := strings.Trim(contact.Name, " ")
	contactFirstName := strings.Trim(contact.FirstName, " ")
	contactLastName := strings.Trim(contact.LastName, " ")

	// Make sure this category is not already in use.
	if action == "create" {

		c := Contact{}

		// IF we have a contact name we validate against that.
		if len(contactName) > 0 {

			if !db.New().Where("ContactsAccountId = ? AND ContactsName = ?", accountId, contactName).First(&c).RecordNotFound() {
				return errors.New("Contact name is already in use.")
			}

			// Double check casing
			if strings.ToLower(contactName) == strings.ToLower(c.Name) {
				return errors.New("Contact name is already in use.")
			}

		} else if (len(contactFirstName) > 0) && (len(contactLastName) > 0) { // Validate first / last

			if !db.New().Where("ContactsAccountId = ? AND ContactsFirstName = ? AND ContactsLastName = ?", accountId, contactFirstName, contactLastName).First(&c).RecordNotFound() {
				return errors.New("Contact first and last name is already in use.")
			}

			// Double check casing
			if strings.ToLower(contactFirstName) == strings.ToLower(c.FirstName) {
				return errors.New("Contact first and last name is already in use.")
			}

			if strings.ToLower(contactLastName) == strings.ToLower(c.LastName) {
				return errors.New("Contact first and last name is already in use.")
			}
		} else {
			return errors.New("A company name or contact first and last name is required.")
		}

	} else if action == "update" {

		c := Contact{}

		if !db.New().Where("ContactsAccountId = ? AND ContactsName = ? AND ContactsFirstName = ? AND ContactsLastName = ?", accountId, contactName, contactFirstName, contactLastName).First(&c).RecordNotFound() {

			// Make sure it is not the same id as the one we are updating
			if c.Id != objId {
				return errors.New("Contact company name, first, and last name is already in use.")
			}
		}

		// Double check casing
		if (c.Id != objId) && (strings.ToLower(contactName) == strings.ToLower(c.Name) && (strings.ToLower(contactFirstName) == strings.ToLower(c.FirstName)) && (strings.ToLower(contactLastName) == strings.ToLower(c.LastName))) {
			return errors.New("Contact company name, first, and last name is already in use.")
		}

	}

	// All good in the hood
	return nil
}

//
// GetContactByAccountAndId - Return a contact by account and id.
//
func (db *DB) GetContactByAccountAndId(accountId uint, conId uint) (Contact, error) {
	l := Contact{}

	// Make query
	if db.New().Where("ContactsAccountId = ? AND ContactsId = ?", accountId, conId).First(&l).RecordNotFound() {
		return Contact{}, errors.New("Contact not found.")
	}

	// Double check the contact has an avatar. This is just to double check.
	db.ConfirmContactAvatar(&l)

	// Add a signed avatar path
	l.AvatarUrl = db.GetSignedFileUrl(l.Avatar)

	// Return result
	return l, nil
}

//
// DeleteContactByAccountAndId - Delete a contact by account and id.
//
func (db *DB) DeleteContactByAccountAndId(accountId uint, contactId uint) error {
	// Make query to see if we have ledger entries with this contact
	if !db.New().Where("LedgerAccountId = ? AND LedgerContactId = ?", accountId, contactId).First(&Ledger{}).RecordNotFound() {
		return errors.New("Can not delete contact. It is in use by a ledger entry.")
	}

	// Make query
	db.New().Where("ContactsAccountId = ? AND ContactsId = ?", accountId, contactId).Delete(Contact{})

	// Return result
	return nil
}

//
// ConfirmContactAvatar will just double check the contact has an avatar
//
func (db *DB) ConfirmContactAvatar(contact *Contact) error {
	// If this is not empty return
	if len(contact.Avatar) > 0 {
		return nil
	}

	// Figure out name
	name := contact.Name

	if len(contact.Name) == 0 {
		name = fmt.Sprintf("%s %s", contact.FirstName, contact.LastName)
	}

	// Guess we do not have an avatar.... create one.
	avatarPath, err := GenerateAndStoreAvatar(contact.AccountId, contact.Id, name, contact.Email)

	if err != nil {
		return err
	}

	// Update DB with the new path of avatar
	contact.Avatar = avatarPath
	db.Save(&contact)

	// Return happy
	return nil
}

// -------------  helper Functions ------------------ //

//
// GenerateAvatarsForAllMissing - run through all accounts and build an avatar for each contact. That does not have one.
// Currently this is called from the CLI.
//
// go run main.go --cmd=contacts-build-missing-avatars
//
func (db *DB) GenerateAvatarsForAllMissing() error {
	// Worker vars.
	jobs := make(chan generateAvatarsWorkerJob, 100000000)
	results := make(chan int, 100000000)

	// Build our our workers
	for w := 1; w <= 5; w++ {
		go db.GenerateAvatarsForAllMissingWoker(jobs, results)
	}

	// Contact var
	c := []Contact{}

	// Run query.
	db.New().Where("ContactsAvatar = ''").Find(&c)

	// Set total count
	count := len(c)

	// Loop through results and make avatar
	for key, row := range c {
		// Figure out name
		name := row.Name

		if len(row.Name) == 0 {
			name = fmt.Sprintf("%s %s", row.FirstName, row.LastName)
		}

		// Test for no name, no contact.
		if len(strings.Trim(name, " ")) == 0 {
			name = "UnKnown Contact"
		}

		// Insert job so worker can run with it.
		jobs <- generateAvatarsWorkerJob{
			name:      name,
			count:     count,
			rowCount:  (key + 1),
			accountId: row.AccountId,
			id:        row.Id,
			email:     row.Email,
		}

		// // Guess we do not have an avatar.... create one.
		// avatarPath, err := GenerateAndStoreAvatar(row.AccountId, row.Id, name, row.Email)
		//
		// if err != nil {
		// 	return err
		// }
		//
		// // Update DB with the new path of avatar
		// row.Avatar = avatarPath
		// db.Save(&row)
		//
		// // Log progress
		// services.Info(fmt.Errorf("Creating avatar for contact AccountId: %d, Id: %d - (%d/%d)", row.AccountId, row.Id, (key + 1), count))
	}

	// Close our jobs channel
	close(jobs)

	// Just pull in the jobs as they are done.
	for i := 1; i <= count; i++ {
		<-results
	}

	// Return happy.
	return nil
}

//
// GenerateAvatarsForAllMissingWoker - Worker to create avatar
//
func (db *DB) GenerateAvatarsForAllMissingWoker(jobs <-chan generateAvatarsWorkerJob, results chan<- int) {
	// Loop through the jobs in the queue.
	for j := range jobs {
		// Guess we do not have an avatar.... create one.
		avatarPath, err := GenerateAndStoreAvatar(j.accountId, j.id, j.name, j.email)

		if err != nil {
			services.Critical(err)
			continue
		}

		// Update DB with the new path of avatar
		contact := Contact{}
		db.New().Find(&contact, j.id)
		contact.Avatar = avatarPath
		db.Save(&contact)

		// Log progress
		services.Info(fmt.Errorf("Creating avatar for contact AccountId: %d, Id: %d - (%d/%d)", j.accountId, j.id, j.rowCount, j.count))

		// Return the job as done.
		results <- j.rowCount
	}
}

//
// GenerateAndStoreAvatar - Generate avatar
//
func GenerateAndStoreAvatar(accountId uint, contactId uint, name string, email string) (string, error) {
	// Wehre we store before upload to S3
	up := ""
	filePath := ""

	// File cache dir.
	cacheDir := fmt.Sprintf("%s/avatars/%d", os.Getenv("CACHE_DIR"), accountId)

	// Make the directory we store this file to
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.MkdirAll(cacheDir, 0755)
	}

	// First we check gravatar
	if len(email) > 0 {
		g, err := gravatar.New()

		if err != nil {
			return "", err
		}

		// Set file file path
		filePath = fmt.Sprintf("%s/%d.jpg", cacheDir, contactId)

		// Get image from gravatar
		g.SetSize(uint(600))
		g.DownloadToDisk(email, filePath)

		// Set upload path
		up = fmt.Sprintf("accounts/%d/avatars/%d.jpg", accountId, contactId)

		// Check the hash of the gravatar to make sure it is not just the default image.
		if files.Md5(filePath) == "f26fffbc0d97cfbe47702676eb7ef799" {
			// Delete uploaded file
			err = os.Remove(filePath)

			if err != nil {
				services.Info(err)
			}

			filePath = ""
		}
	}

	// Build default Avatar
	if len(filePath) == 0 {
		filePath = fmt.Sprintf("%s/%d.png", cacheDir, contactId)
		err := avatar.ToDisk(name, filePath)

		if err != nil {
			return "", err
		}

		// Set upload path
		up = fmt.Sprintf("accounts/%d/avatars/%d.png", accountId, contactId)
	}

	// Upload file to our S3 store
	err := object.UploadObject(filePath, up)

	if err != nil {
		return "", fmt.Errorf("generateAndStoreAvatar: ContactId: %d, AccountId: %d Error: %s", contactId, accountId, err.Error())
	}

	// Delete uploaded file
	err = os.Remove(filePath)

	if err != nil {
		services.Info(err)
	}

	// Return happy
	return up, nil
}

/* End File */
