//
// Date: 10/26/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"app.skyclerk.com/backend/library/realip"
	"app.skyclerk.com/backend/services"
	"github.com/gin-gonic/gin"
)

//
// Post back to setup a forgot password request. - Step #1 (send email request to reset)
//
func (t *Controller) DoForgotPassword(c *gin.Context) {
	// Set response
	if os.Getenv("APP_ENV") == "local" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// Decode json passed in
	decoder := json.NewDecoder(c.Request.Body)

	type ForgotPost struct {
		Email string
	}

	var post ForgotPost

	err := decoder.Decode(&post)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."})
		return
	}

	defer c.Request.Body.Close()

	// Request a reset password request.
	err = t.db.DoResetPassword(post.Email, realip.RealIP(c.Request))

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, we could not find your account."})
		return
	}

	// Return success json.
	c.JSON(http.StatusNoContent, nil)
}

//
// Rest the password after they clicked on the email - Step #2
//
func (t *Controller) DoResetPassword(c *gin.Context) {

	// Set response
	if os.Getenv("APP_ENV") == "local" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// Decode json passed in
	decoder := json.NewDecoder(c.Request.Body)

	type ResetPost struct {
		Hash     string
		Password string
	}

	var post ResetPost

	err := decoder.Decode(&post)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."})
		return
	}

	defer c.Request.Body.Close()

	// Get the user based on the hash we passed in.
	user, err := t.db.GetUserFromToken(post.Hash)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, It seems your reset token has expired."})
		return
	}

	// Now that we know the user lets make sure the password that was posted in was at least 6 chars.
	err = t.db.ValidatePassword(post.Password)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a password at least 6 chars long."})
		return
	}

	// Now that we know who the user is lets reset the users password.
	err = t.db.ResetUserPassword(user.Id, post.Password)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."})
		return
	}

	// Lastly delete the reset password hash.
	err = t.db.DeleteForgotPasswordByToken(post.Hash)

	if err != nil {
		services.Info(err)
	}

	// Return success json.
	c.JSON(http.StatusNoContent, nil)
}

/* End File */
