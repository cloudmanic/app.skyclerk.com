//
// Date: 4/14/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"

	"app.skyclerk.com/backend/library/helpers"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"
)

//
// GetMe returns all the data related to a user.
//
func (t *Controller) GetMe(c *gin.Context) {
	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(int)

	// Get the full user
	user, err := t.db.GetUserById(uint(userId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}

	// Return happy JSON
	c.JSON(200, user)
}

//
// ChangePassword for the logged in user
//
func (t *Controller) ChangePassword(c *gin.Context) {
	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(int)

	// Get the full user
	user, err := t.db.GetUserById(uint(userId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}

	// Get JSON passed in.
	body, _ := ioutil.ReadAll(c.Request.Body)
	current := gjson.Get(string(body), "current").String()
	password := gjson.Get(string(body), "password").String()
	confirm := gjson.Get(string(body), "confirm").String()

	// Validate the passwords match.
	if password != confirm {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match."})
		return
	}

	// Validate password
	verr := t.db.ValidatePassword(password)

	if verr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": verr.Error()})
		return
	}

	// Validate current password.
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(current))

	if err != nil {
		// Check MD5 login
		if (len(user.Md5Salt) > 0) && (len(user.Md5Password) > 0) {
			passMd5 := helpers.GetMd5(current + user.Md5Salt)

			if passMd5 != user.Md5Password {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Your current password is not correct."})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Your current password is not correct."})
			return
		}
	}

	// Create new hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to change your password. Please email help@skyclerk.com."})
		return
	}

	// Change password
	user.Md5Password = ""
	user.Md5Salt = ""
	user.Password = string(hash)
	t.db.New().Save(&user)

	// Return happy JSON
	c.JSON(204, nil)
}

//
// UpdateMe will update the user profile.
//
func (t *Controller) UpdateMe(c *gin.Context) {
	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(int)

	// Get the full user
	user, err := t.db.GetUserById(uint(userId))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}

	// Get JSON passed in.
	body, _ := ioutil.ReadAll(c.Request.Body)
	first := gjson.Get(string(body), "first_name").String()
	last := gjson.Get(string(body), "last_name").String()
	email := gjson.Get(string(body), "email").String()

	// Make sure we have all the fields.
	if len(first) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "First name field is required."})
		return
	}

	if len(last) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Last name field is required."})
		return
	}

	if len(email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email field is required."})
		return
	}

	// Lets validate the email address
	if err := t.db.ValidateEmailAddress(email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// See if we already have this user.
	if user.Email != email {
		_, err := t.db.GetUserByEmail(strings.Trim(email, " "))

		// Meaning we found the user so we can't use the email again.
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use."})
			return
		}
	}

	// Change password
	user.FirstName = strings.Trim(first, " ")
	user.LastName = strings.Trim(last, " ")
	user.Email = strings.Trim(email, " ")
	t.db.New().Save(&user)

	// Return happy JSON
	c.JSON(204, nil)
}
