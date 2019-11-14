package models

import (
	"errors"
	"os"
	"time"

	"app.skyclerk.com/backend/library/email"
	"app.skyclerk.com/backend/services"
)

type ForgotPassword struct {
	Id        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    uint   `sql:"not null;index:UserId"`
	Token     string `sql:"not null"`
	IpAddress string `sql:"not null"`
}

//
// Return the user based on the hash passed.
//
func (t *DB) GetUserFromToken(token string) (User, error) {

	var u User
	var f ForgotPassword

	// Get the record based on the token we passed in.
	if t.Where("Token = ?", token).First(&f).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Get the user from based on the user id we got.
	u, err := t.GetUserById(f.UserId)

	if err != nil {
		return u, err
	}

	// Return the user.
	return u, nil

}

//
// Delete record by hash.
//
func (t *DB) DeleteForgotPasswordByToken(token string) error {

	var f ForgotPassword

	// Get the record based on the token we passed in.
	if t.Where("Token = ?", token).First(&f).RecordNotFound() {
		return errors.New("Record not found")
	}

	// Delete record
	if err := t.Delete(&f).Error; err != nil {
		return err
	}

	// Return success
	return nil

}

//
// DoResetPassword Reset the user's password and send an email telling them next steps.
//
func (t *DB) DoResetPassword(user_email string, ip string) error {
	// Make sure this is a real email address
	user, err := t.GetUserByEmail(user_email)

	if err != nil {
		return errors.New("Sorry, we were unable to find our account.")
	}

	// Generate "hash" to store for the reset token
	hash, err := t.GenerateRandomString(30)

	if err != nil {
		services.Critical(errors.New(err.Error() + "DoResetPassword - Unable to create hash (GenerateRandomString)"))
		return err
	}

	// Store the new reset password hash.
	rsp := ForgotPassword{UserId: user.Id, Token: hash, IpAddress: ip}
	t.Create(&rsp)

	// Log user creation.
	services.InfoMsg("DoResetPassword - Reset password token for " + user.Email)

	// Build the url to reset the password.
	var url = os.Getenv("SITE_URL") + "/reset-password?hash=" + hash

	// Build attachments.
	attachments := []string{}

	// Send email to user asking them to come to the site and reset the password.
	err = email.Send(
		user.Email,
		"Reset Your Password",
		t.GetForgotPasswordStepOneEmailHtml(user.FirstName, user.Email, url),
		attachments,
	)

	if err != nil {
		return err
	}

	// Everything went as planned.
	return nil

}

// ------------------- Template Emails ------------------------- //

//
// Get the html version of the forgot password email step #1
//
func (t *DB) GetForgotPasswordStepOneEmailHtml(name string, email string, url string) string {

	return string(`
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office"><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta name="Generator" content="Made with Mail Designer from equinux">
	<meta name="Viewport" content="width=device-width, initial-scale=1.0">
	<style type="text/css" id="Mail Designer General Style Sheet">
		a { word-break: break-word; }
		a img { border:none; }
		img { outline:none; text-decoration:none; -ms-interpolation-mode: bicubic; }
		body { width: 100% !important; -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; }
		.ExternalClass { width: 100%; }
		.ExternalClass, .ExternalClass p, .ExternalClass span, .ExternalClass font, .ExternalClass td, .ExternalClass div { line-height: 100%; }
		#page-wrap { margin: 0; padding: 0; width: 100% !important; line-height: 100% !important; }
		#outlook a { padding: 0; }
		.preheader { display:none !important; }
		a[x-apple-data-detectors] { color: inherit !important; text-decoration: none !important; font-size: inherit !important; font-family: inherit !important; font-weight: inherit !important; line-height: inherit !important; }
		.a5q { display: none !important; }
		.Apple-web-attachment { vertical-align: initial !important; }
		.Apple-edge-to-edge-visual-media { margin: initial !important; max-width: initial !important; width: 100%; }
	</style>
	<style type="text/css" id="Mail Designer Mobile Style Sheet">
		@media only screen and (max-width: 580px) {
			table.email-body-wrap {
				width: 320px !important;
			}
			td.page-bg-show-thru {
				display: none !important;
			}
			table.layout-block-wrapping-table {
				width: 320px !important;
			}
			table.mso-fixed-width-wrapping-table {
				width: 320px !important;
			}
			.layout-block-full-width {
				width: 320px !important;
			}
			table.layout-block-column, table.layout-block-padded-column {
				width: 100% !important;
			}
			table.layout-block-box-padding {
				width: 100% !important;
				padding: 5px !important;
			}
			table.layout-block-horizontal-spacer {
				display: none !important;
			}
			tr.layout-block-vertical-spacer {
			   display: block !important;
			   height: 8px !important;
			}
			td.container-padding {
				display: none !important;
			}

			table {
				min-width: initial !important;
			}
			td {
				min-width: initial !important;
			}

			.desktop-only { display: none !important; }
			.mobile-only { display: block !important; }

			.hide {
				max-height: none !important;
				display: block !important;
				overflow: visible !important;
			}

			.EQ-00 { width: 320px !important; }
			.EQ-01 { width: 320px !important; height: 32px !important; }

			.EQ-02 { width: 7px !important; }
			.EQ-03 { width: 16px !important; }
			.EQ-04 { width: 297px !important; }
			.EQ-05 { height:12px !important; } /* vertical spacer */

			.EQ-06 { width: 12px !important; }
			.EQ-07 { width: 6px !important; }
			.EQ-08 { width: 302px !important; }

			.EQ-09 { width: 6px !important; }
			.EQ-0A { width: 6px !important; }
			.EQ-0B { width: 308px !important; }
			.EQ-0C { width: 298px !important; height: 46px !important; }

			.EQ-0D { width: 12px !important; }
			.EQ-0E { width: 6px !important; }
			.EQ-0F { width: 302px !important; }

			.EQ-10 { width: 6px !important; }
			.EQ-11 { width: 6px !important; }
			.EQ-12 { width: 308px !important; }
			.EQ-13 { height:12px !important; } /* vertical spacer */
		}
	</style>
	<!--[if gte mso 9]>
	<style type="text/css" id="Mail Designer Outlook Style Sheet">
		table.layout-block-horizontal-spacer {
		    display: none !important;
		}
		table {
		    border-collapse:collapse;
		    mso-table-lspace:0pt;
		    mso-table-rspace:0pt;
		    mso-table-bspace:0pt;
		    mso-table-tspace:0pt;
		    mso-padding-alt:0;
		    mso-table-top:0;
		    mso-table-wrap:around;
		}
		td {
		    border-collapse:collapse;
		    mso-cellspacing:0;
		}
	</style>
	<xml>
		<o:OfficeDocumentSettings>
			<o:AllowPNG/>
			<o:PixelsPerInch>96</o:PixelsPerInch>
		</o:OfficeDocumentSettings>
	</xml>
	<![endif]-->
<link href="https://fonts.googleapis.com/css?family=Droid+Sans:700,regular" rel="stylesheet" type="text/css" class="EQWebFont"><link href="https://fonts.googleapis.com/css?family=Roboto:regular,700" rel="stylesheet" type="text/css" class="EQWebFont"><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head>
<body style="margin-top: 0px; margin-right: 0px; margin-bottom: 0px; margin-left: 0px; padding-top: 0px; padding-right: 0px; padding-bottom: 0px; padding-left: 0px; " background="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1"><!--[if gte mso 9]>
<v:background xmlns:v="urn:schemas-microsoft-com:vml" fill="t">
<v:fill type="tile" src="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" />
</v:background>
<![endif]-->

<table width="100%" cellspacing="0" cellpadding="0" id="page-wrap" align="center" background="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1">
<tbody><tr>
	<td>

<table class="email-body-wrap" width="610" cellspacing="0" cellpadding="0" id="email-body" align="center">
<tbody><tr>
	<td width="30" class="page-bg-show-thru">&nbsp;<!--Left page bg show-thru --></td>
	<td width="550" id="page-body">

		<!--Begin of layout container -->
		<div id="eqLayoutContainer">

			<div width="100%">
				<!--[if !mso 15]><!--><table width="550" cellspacing="0" cellpadding="0" class="mso-fixed-width-wrapping-table" style="mso-hide:all; min-width: 550px; ">
					<tbody><tr>
						<td valign="top" class="layout-block-full-width" width="550" style="min-width: 550px; ">
							<table cellspacing="0" cellpadding="0" class="layout-block-full-width" width="550" style="min-width: 550px; ">
								<tbody><tr>
									<td width="550" style="min-width: 550px; ">
										<div class="layout-block-image">
											<a href="https://skyclerk.com"><img width="550" height="55" alt="" src="https://cdn.skyclerk.com/email/forgot-password-step-1/image-1.png" border="0" style="display: block; width: 550px; height: 55px; " class="EQ-01"></a>
										</div>
									</td>
								</tr>
							</tbody></table>
						</td>
					</tr>
				</tbody></table>
			<!--<![endif]--><!--[if gte mso 9]><table cellspacing="0" cellpadding="0" class="mso-fixed-width-wrapping-table" width="550">
<tr><td valign="top" class="layout-block-full-width"><div class="layout-block-image"><a href="https://skyclerk.com"><img height="55" alt="" src="https://cdn.skyclerk.com/email/forgot-password-step-1/image-1.png" border="0" style="display: block;"  class="img-EQMST-414F061F-0A35-4461-A492-2F9B6543416D" width="550"></img><div style="width:0px;height:0px;max-height:0;max-width:0;overflow:hidden;display:none;visibility:hidden;mso-hide:all;"></div></a></div></td>
</tr></table><![endif]--></div>

			<div width="100%">
				<!--[if !mso 15]><!--><table width="550" cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" style="mso-hide:all; min-width: 550px; ">
					<tbody><tr>
						<td width="12" class="EQ-02" style="font-size: 1px; min-width: 12px; ">&nbsp;</td>
						<td height="20" background="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" class="EQ-04" bgcolor="#ffffff" width="511" style="min-width: 511px; ">
							<div class="spacer"></div>
						</td>
						<td width="27" class="EQ-03" style="font-size: 1px; min-width: 27px; ">&nbsp;</td>
					</tr>
				</tbody></table>
			<!--<![endif]--><!--[if gte mso 9]><table cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" width="550">
<tr><td width="12" class="layout-block-padding-left" style="font-size:0">&nbsp; </td><td class="layout-block-content-cell" style="font-size:0;" width="511" height="20"><v:rect style="width:511px;height:20px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" color="#ffffff"></v:fill><div class="spacer"></div></v:rect></td><td width="27" class="layout-block-padding-right" style="font-size:0">&nbsp; </td>
</tr></table><![endif]--></div>

			<div width="100%">
				<!--[if !mso 15]><!--><table width="550" cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" style="mso-hide:all; min-width: 550px; ">
					<tbody><tr>
						<td width="20" class="EQ-06" style="min-width: 20px; ">&nbsp;</td>
						<td width="520" class="EQ-08" style="min-width: 520px; ">
							<table cellspacing="0" cellpadding="0" align="left" class="layout-block-column" width="520" style="min-width: 520px; ">
								<tbody><tr>
									<td width="500" valign="top" align="left" background="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" style="padding-left: 10px; padding-right: 10px; min-width: 500px; " bgcolor="#ffffff">
										<table cellspacing="0" cellpadding="0" class="layout-block-box-padding" width="500" style="min-width: 500px; ">
											<tbody><tr>
												<td align="left" class="layout-block-column" width="500" style="min-width: 500px; ">
													<div class="text" style="font-size: 16px; font-family: 'Lucida Grande'; line-height: 1.2; "><font face="Roboto, Times New Roman, sans-serif" style="line-height: 1.2; font-size: 17px; ">Hi ` + name + `,</font><div style="line-height: 1.2; "><font face="Roboto, Times New Roman, sans-serif" style="font-size: 17px; "><br></font></div><div style="line-height: 1.2; "><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; font-size: 17px; "><font face="Roboto, helvetica, arial, sans-serif">Can't remember your password? Don't worry about it — it happens.</font></span></div><div style="line-height: 1.2; "><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; font-size: 17px; "><font face="Roboto, helvetica, arial, sans-serif"><br></font></span></div><div style="line-height: 1.2; "><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; font-size: 17px; "><font face="Roboto, helvetica, arial, sans-serif"><span style="font-variant-ligatures: normal; letter-spacing: -0.1px; ">So you know, your username is: </span><span style="font-variant-ligatures: normal; font-variant-numeric: inherit; font-stretch: inherit; line-height: inherit; letter-spacing: -0.1px; "><b>` + email + `</b></span></font></span></div></div>
												</td>
											</tr>
										</tbody></table>
									</td>
								</tr>
							</tbody></table>
						</td>
						<td width="10" class="EQ-07" style="min-width: 10px; ">&nbsp;</td>
					</tr>
				</tbody></table>
			<!--<![endif]--><!--[if gte mso 9]><table cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" width="550">
<tr><td width="20" class="layout-block-padding-left">&nbsp; </td><td width="520" class="layout-block-content-cell" valign="top" align="left"><v:rect style="width:520px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" color="#ffffff"></v:fill><v:textbox style="mso-fit-shape-to-text:true" inset="0,0,0,0"><div><div style="font-size:0"><table cellspacing="0" cellpadding="0" class="layout-block-box-padding">
<tr><td style="font-size:1px;" width="10">&nbsp; </td><td align="left" class="layout-block-column" width="500"><div class="text" style="font-size: 16px; font-family: sans-serif; line-height: 120%;"><font face="Times New Roman, sans-serif" style="line-height: 120%; font-size: 17px;">Hi John,</font><div style="line-height: 120%;"><font face="Times New Roman, sans-serif" style="font-size: 17px;"><br></font></div><div style="line-height: 120%;"><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; font-size: 17px;"><font face="sans-serif">Can't remember your password? Don't worry about it — it happens.</font></span></div><div style="line-height: 120%;"><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; font-size: 17px;"><font face="sans-serif"><br></font></span></div><div style="line-height: 120%;"><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; font-size: 17px;"><font face="sans-serif"><span style="font-variant-ligatures: normal; letter-spacing: -0.1px;">So you know, your username is: </span><span style="font-variant-ligatures: normal; font-variant-numeric: inherit; font-stretch: inherit; line-height: inherit; letter-spacing: -0.1px;"><b>` + email + `</b></span></font></span></div></div></td><td style="font-size:1px;" width="10">&nbsp; </td>
</tr></table></div></div></v:textbox></v:rect></td><td width="10" class="layout-block-padding-right">&nbsp; </td>
</tr></table><![endif]--></div>

			<div width="100%">
				<!--[if !mso 15]><!--><table width="550" cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" style="mso-hide:all; min-width: 550px; ">
					<tbody><tr>
						<td width="10" class="EQ-09" style="min-width: 10px; ">&nbsp;</td>
						<td background="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" class="EQ-0B" bgcolor="#ffffff" width="530" style="min-width: 530px; ">
							<table cellspacing="0" cellpadding="0" align="left" style="padding-left: 10px; padding-right: 10px; min-width: 530px; " class="layout-block-box-padding" width="530">
								<tbody><tr>
									<td valign="top" class="layout-block-padded-column" align="center" width="510" style="min-width: 510px; ">
										<div class="layout-block-image">
											<a href="` + url + `"><img width="510" height="79" alt="Skyclerk Reset Password" src="https://cdn.skyclerk.com/email/forgot-password-step-1/image-2.png" border="0" style="display: block; width: 510px; height: 79px; " class="EQ-0C"></a>
										</div>
									</td>
								</tr>
							</tbody></table>
						</td>
						<td width="10" class="EQ-0A" style="min-width: 10px; ">&nbsp;</td>
					</tr>
				</tbody></table>
			<!--<![endif]--><!--[if gte mso 9]><table cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" width="550">
<tr><td width="10" class="layout-block-padding-left">&nbsp; </td><td class="layout-block-content-cell" width="530"><v:rect style="width:530px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" color="#ffffff"></v:fill><v:textbox style="mso-fit-shape-to-text:true" inset="0,0,0,0"><div><div style="font-size:0"><table cellspacing="0" cellpadding="0" align="left" style="padding-left:10px; padding-right:10px;" class="layout-block-box-padding">
<tr><td valign="top" class="layout-block-padded-column" align="center"><div class="layout-block-image"><a href="` + url + `"><img height="79" alt="Skyclerk Reset Password" src="https://cdn.skyclerk.com/email/forgot-password-step-1/image-2.png" border="0"  style="display: block;" class="img-EQMST-CCAE6821-F243-43A3-84D0-AFDB050CF64E" width="510"></img><div style="width:0px;height:0px;max-height:0;max-width:0;overflow:hidden;display:none;visibility:hidden;mso-hide:all;"></div></a></div></td>
</tr></table></div></div></v:textbox></v:rect></td><td width="10" class="layout-block-padding-right">&nbsp; </td>
</tr></table><![endif]--></div>

			<div width="100%">
				<!--[if !mso 15]><!--><table width="550" cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" style="mso-hide:all; min-width: 550px; ">
					<tbody><tr>
						<td width="20" class="EQ-0D" style="min-width: 20px; ">&nbsp;</td>
						<td width="520" class="EQ-0F" style="min-width: 520px; ">
							<table cellspacing="0" cellpadding="0" align="left" class="layout-block-column" width="520" style="min-width: 520px; ">
								<tbody><tr>
									<td width="500" valign="top" align="left" background="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" style="padding-left: 10px; padding-right: 10px; min-width: 500px; " bgcolor="#ffffff">
										<table cellspacing="0" cellpadding="0" class="layout-block-box-padding" width="500" style="min-width: 500px; ">
											<tbody><tr>
												<td align="left" class="layout-block-column" width="500" style="min-width: 500px; ">
													<div class="text" style="font-size: 16px; font-family: 'Lucida Grande'; "><div style="text-align: left; "><span style="font-size: 17px; "><font face="Roboto, helvetica, arial, sans-serif"><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; text-align: center; widows: 2; ">Have questions? Need help? </span><a href="https://skyclerk.com/support" style="color: rgb(22, 108, 197); ">Contact our support team</a><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; text-align: center; widows: 2; "> and&nbsp;</span></font><font face="Roboto, helvetica, arial, sans-serif" style="text-align: center; "><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; ">we'll get back to you very quickly</span></font><span style="text-align: center; color: rgb(31, 31, 31); font-family: 'Helvetica Neue', helvetica, arial, sans-serif; font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; ">.</span></span></div><div style="text-align: left; "><span style="text-align: center; color: rgb(31, 31, 31); font-family: 'Helvetica Neue', helvetica, arial, sans-serif; font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; font-size: 17px; "><br></span></div><div style="text-align: left; "><div style="line-height: 1.2; "><font face="Roboto, Times New Roman, sans-serif" style="font-size: 17px; ">Thanks!</font></div><div style="line-height: 1.2; "><font face="Roboto, Times New Roman, sans-serif" style="font-size: 17px; ">- Your Skyclerk Team</font></div><div><br></div></div></div>
												</td>
											</tr>
										</tbody></table>
									</td>
								</tr>
							</tbody></table>
						</td>
						<td width="10" class="EQ-0E" style="min-width: 10px; ">&nbsp;</td>
					</tr>
				</tbody></table>
			<!--<![endif]--><!--[if gte mso 9]><table cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" width="550">
<tr><td width="20" class="layout-block-padding-left">&nbsp; </td><td width="520" class="layout-block-content-cell" valign="top" align="left"><v:rect style="width:520px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" color="#ffffff"></v:fill><v:textbox style="mso-fit-shape-to-text:true" inset="0,0,0,0"><div><div style="font-size:0"><table cellspacing="0" cellpadding="0" class="layout-block-box-padding">
<tr><td style="font-size:1px;" width="10">&nbsp; </td><td align="left" class="layout-block-column" width="500"><div class="text" style="font-size: 16px; font-family: sans-serif;"><div style="text-align: left;"><span style="font-size: 17px;"><font face="sans-serif"><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; text-align: center; widows: 2;">Have questions? Need help? </span><a href="https://skyclerk.com/support" style="color: rgb(22, 108, 197);">Contact our support team</a><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; text-align: center; widows: 2;"> and </span></font><font face="sans-serif" style="text-align: center;"><span style="color: rgb(31, 31, 31); font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2;">we'll get back to you very quickly</span></font><span style="text-align: center; color: rgb(31, 31, 31); font-family: sans-serif; font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2;">.</span></span></div><div style="text-align: left;"><span style="text-align: center; color: rgb(31, 31, 31); font-family: sans-serif; font-variant-ligatures: normal; letter-spacing: -0.1px; orphans: 2; widows: 2; font-size: 17px;"><br></span></div><div style="text-align: left;"><div style="line-height: 120%;"><font face="Times New Roman, sans-serif" style="font-size: 17px;">Thanks!</font></div><div style="line-height: 120%;"><font face="Times New Roman, sans-serif" style="font-size: 17px;">- Your Skyclerk Team</font></div><div><br></div></div></div></td><td style="font-size:1px;" width="10">&nbsp; </td>
</tr></table></div></div></v:textbox></v:rect></td><td width="10" class="layout-block-padding-right">&nbsp; </td>
</tr></table><![endif]--></div>

			<div width="100%">
				<!--[if !mso 15]><!--><table width="550" cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" style="mso-hide:all; min-width: 550px; ">
					<tbody><tr>
						<td width="10" class="EQ-10" style="font-size: 1px; min-width: 10px; ">&nbsp;</td>
						<td height="20" background="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" class="EQ-12" bgcolor="#ffffff" width="530" style="min-width: 530px; ">
							<div class="spacer"></div>
						</td>
						<td width="10" class="EQ-11" style="font-size: 1px; min-width: 10px; ">&nbsp;</td>
					</tr>
				</tbody></table>
			<!--<![endif]--><!--[if gte mso 9]><table cellspacing="0" cellpadding="0" class="layout-block-wrapping-table" width="550">
<tr><td width="10" class="layout-block-padding-left" style="font-size:0">&nbsp; </td><td class="layout-block-content-cell" style="font-size:0;" width="530" height="20"><v:rect style="width:530px;height:20px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/email/forgot-password-step-1/box-bg.jpg?v=1" color="#ffffff"></v:fill><div class="spacer"></div></v:rect></td><td width="10" class="layout-block-padding-right" style="font-size:0">&nbsp; </td>
</tr></table><![endif]--></div>

		</div>
		<!--End of layout container -->

	</td>
	<td width="30" class="page-bg-show-thru">&nbsp;<!--Right page bg show-thru --></td>
</tr>
</tbody></table><!--email-body -->

	</td>
</tr>
</tbody></table><!--page-wrap -->


</body></html>
  `)

}

/* End File */
