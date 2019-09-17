//
// Date: 2019-07-02
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package emails

import (
	"app.skyclerk.com/backend/models"
)

//
// GetInviteNewUserHTML will set html
//
func GetInviteNewUserHTML(name string, accountName string, url string, invite models.Invite) string {
	return `

	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">

	<head>
		<!-- Required for Yahoo Mail app -->
	</head>

	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<meta name="Viewport" content="width=device-width, initial-scale=1.0">
		<style type="text/css" id="Mail Designer General Style Sheet">
			a {
				word-break: break-word;
			}

			a img {
				border: none;
			}

			img {
				outline: none;
				text-decoration: none;
				-ms-interpolation-mode: bicubic;
			}

			body {
				width: 100% !important;
				-webkit-text-size-adjust: 100%;
				-ms-text-size-adjust: 100%;
			}

			.ExternalClass {
				width: 100%;
			}

			.ExternalClass,
			.ExternalClass p,
			.ExternalClass span,
			.ExternalClass font,
			.ExternalClass td,
			.ExternalClass div {
				line-height: 100%;
			}

			#page-wrap {
				margin: 0;
				padding: 0;
				width: 100% !important;
				line-height: 100% !important;
			}

			#outlook a {
				padding: 0;
			}

			.preheader {
				display: none !important;
			}

			a[x-apple-data-detectors] {
				color: inherit !important;
				text-decoration: none !important;
				font-size: inherit !important;
				font-family: inherit !important;
				font-weight: inherit !important;
				line-height: inherit !important;
			}

			.a5q {
				display: none !important;
			}

			.Apple-web-attachment {
				vertical-align: initial !important;
			}

			.Apple-edge-to-edge-visual-media {
				margin: initial !important;
				max-width: initial !important;
				width: 100%;
			}

			ul {
				margin-top: 0;
				margin-bottom: 0;
				padding-top: 0;
				padding-bottom: 0;
			}

			ol {
				margin-top: 0;
				margin-bottom: 0;
				padding-top: 0;
				padding-bottom: 0;
			}

		</style>
		<style type="text/css" id="Mail Designer Mobile Style Sheet">
			@media only screen and (max-width: 580px) {
				table.EQ-00 {
					width: 320px !important;
				}

				td.EQ-01 {
					display: none !important;
				}

				.EQ-04 {
					width: 320px !important;
				}

				table.EQ-05,
				table.EQ-06 {
					width: 100% !important;
				}

				table.EQ-07 {
					width: 100% !important;
					padding: 5px !important;
				}

				table.layout-block-horizontal-spacer {
					display: none !important;
				}

				tr.EQ-08 {
					display: block !important;
					height: 8px !important;
				}

				table {
					min-width: initial !important;
				}

				td {
					min-width: initial !important;
				}

				.EQ-10 {
					display: none !important;
				}

				.mobile-only {
					display: block !important;
				}

				.EQ-11 {
					max-height: none !important;
					display: block !important;
					overflow: visible !important;
				}

				.layout-block-table-desktop {
					display: none !important;
				}

				.layout-block-table-mobile {
					width: 100% !important;
					display: block !important;
				}

				.md-table-spacer {
					height: 50px;
				}

				#eqLayoutContainer {}

				table.EQ-12 {
					padding-top: 0 !important;
				}

				table.EQ-13 {
					padding-right: 0 !important;
				}

				table.EQ-14 {
					padding-bottom: 0 !important;
				}

				table.EQ-15 {
					padding-left: 0 !important;
				}

				.EQ-16 {
					width: 320px !important;
				}

				.EQ-17 {
					width: 320px !important;
					height: 51px !important;
				}

				.EQ-18 {
					width: 7px !important;
				}

				.EQ-19 {
					width: 16px !important;
				}

				.EQ-20 {
					width: 297px !important;
				}

				.EQ-21 {
					height: 12px !important;
				}

				.EQ-22 {
					width: 12px !important;
				}

				.EQ-23 {
					width: 6px !important;
				}

				.EQ-24 {
					width: 302px !important;
				}

				.EQ-28 {
					width: 267px !important;
					height: 47px !important;
				}

				.EQ-34 {
					width: 308px !important;
				}
			}

		</style>
		<!--[if !mso 15]><!-->
		<style type="text/css" id="Outlook hidden">
			#page-wrap {
				background-color: rgb(255, 255, 255);
			}

		</style>
		<!--<![endif]-->
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
		<link href="https://fonts.googleapis.com/css?family=Droid+Sans:700,regular" rel="stylesheet" type="text/css" class="EQWebFont">
		<link href="https://fonts.googleapis.com/css?family=Roboto:regular,700,italic" rel="stylesheet" type="text/css" class="EQWebFont">
		<style type="text/css" id="md365-mobile-modified">

		</style>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	</head>

	<body style="margin-top: 0px; margin-right: 0px; margin-bottom: 0px; margin-left: 0px; padding-top: 0px; padding-right: 0px; padding-bottom: 0px; padding-left: 0px;" background="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg">
		<!--[if gte mso 9]>
	<v:background xmlns:v="urn:schemas-microsoft-com:vml" fill="t">
	<v:fill type="tile" src="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg" />
	</v:background>
	<![endif]-->

		<table width="100%" cellspacing="0" cellpadding="0" id="page-wrap" align="center" background="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg">
			<tbody>
				<tr>
					<td>

						<table class="EQ-00" width="610" cellspacing="0" cellpadding="0" id="email-body" align="center">
							<tbody>
								<tr>
									<td width="30" class="EQ-01">&nbsp;
										<!--Left page bg show-thru -->
									</td>
									<td width="550" id="page-body">

										<!--Begin of layout container -->
										<div id="eqLayoutContainer">

											<div width="100%">
												<!--[if !mso 15]><!-->
												<table width="550" cellspacing="0" cellpadding="0" class="EQ-00" style="mso-hide:all;">
													<tbody>
														<tr>
															<td valign="top" class="EQ-04" width="550">
																<table cellspacing="0" cellpadding="0" class="EQ-04" width="550">
																	<tbody>
																		<tr>
																			<td width="550">
																				<div class="layout-block-image">
																					<a href="https://skyclerk.com" target="_blank"><img width="550" height="88" alt="" src="https://cdn.skyclerk.com/emails/invite-new-user/image-1.png" border="0" style="display: block; width: 550px; height: 88px;" class="EQ-17"></a>
																				</div>
																			</td>
																		</tr>
																	</tbody>
																</table>
															</td>
														</tr>
													</tbody>
												</table>
												<!--<![endif]-->
												<!--[if gte mso 9]><div style="display:none;"><table cellspacing="0" cellpadding="0" class="EQ-00" width="550">
	<tr><td valign="top" class="EQ-04"><div class="layout-block-image"><a href="https://skyclerk.com"><img height="88" alt="" src="https://cdn.skyclerk.com/emails/invite-new-user/image-1.png" border="0" style="display: block;"  class="EQ-17" width="550"></img><div style="width:0px;height:0px;max-height:0;max-width:0;overflow:hidden;display:none;visibility:hidden;mso-hide:all;"></div></a></div></td>
	</tr></table></div><![endif]-->
											</div>

											<div width="100%">
												<!--[if !mso 15]><!-->
												<table width="550" cellspacing="0" cellpadding="0" class="EQ-00" style="mso-hide:all;">
													<tbody>
														<tr>
															<td width="12" class="EQ-18" style="font-size:1px;">&nbsp;</td>
															<td height="20" background="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg" class="EQ-20" bgcolor="#ffffff" width="511">
																<div class="spacer"></div>
															</td>
															<td width="27" class="EQ-19" style="font-size:1px;">&nbsp;</td>
														</tr>
													</tbody>
												</table>
												<!--<![endif]-->
												<!--[if gte mso 9]><div style="display:none;"><table cellspacing="0" cellpadding="0" class="EQ-00" width="550">
	<tr><td width="12" class="layout-block-padding-left" style="font-size:0">&nbsp; </td><td class="layout-block-content-cell" style="font-size:0;" width="511" height="20"><v:rect style="width:511px;height:20px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg" color="#ffffff"></v:fill><div class="spacer"></div></v:rect></td><td width="27" class="layout-block-padding-right" style="font-size:0">&nbsp; </td>
	</tr></table></div><![endif]-->
											</div>

											<div width="100%">
												<!--[if !mso 15]><!-->
												<table width="550" cellspacing="0" cellpadding="0" class="EQ-00" style="mso-hide:all;">
													<tbody>
														<tr>
															<td width="20" class="EQ-22">&nbsp;</td>
															<td width="520" class="EQ-24">
																<table cellspacing="0" cellpadding="0" align="left" class="EQ-05">
																	<tbody>
																		<tr>
																			<td width="500" valign="top" align="left" background="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg" style="padding-left: 10px; padding-right: 10px;" bgcolor="#ffffff">
																				<table cellspacing="0" cellpadding="0" class="EQ-07">
																					<tbody>
																						<tr>
																							<td align="left" class="EQ-05" width="500">
																								<div class="text" style="font-size: 16px; font-family: 'Lucida Grande'; line-height: 1.2;">
																									<font face="Roboto, Times New Roman, sans-serif" style="line-height: 1.2;">Hi ` + invite.FirstName + `,</font>
																									<div style="line-height: 1.2;"><b style="font-family: Times; font-size: 16px;"><br></b></div>
																									<div style="line-height: 1.2;">
																										<font face="Roboto, Times New Roman, sans-serif"><b style="font-size: 16px;">` + name + `</b><span style="font-size: 16px;">&nbsp;invited you to the&nbsp;</span><b style="font-size: 16px;">` + accountName + `</b><span style="font-size: 16px;"> account.</span></font><span style="font-family: Roboto, &quot;Times New Roman&quot;, sans-serif;">&nbsp;</span>
																									</div>
																									<div style="line-height: 1.2;"><span style="font-family: Roboto, &quot;Times New Roman&quot;, sans-serif;"><br></span></div>
																									<div style="line-height: 1.2;"><span style="font-family: Roboto, helvetica, &quot;Times New Roman&quot;, sans-serif; font-size: 15px; font-style: italic;">` + invite.Message + `</span></div>
																									<div style="line-height: 1.2;"><br></div>
																								</div>
																							</td>
																						</tr>
																					</tbody>
																				</table>
																			</td>
																		</tr>
																	</tbody>
																</table>
															</td>
															<td width="10" class="EQ-23">&nbsp;</td>
														</tr>
													</tbody>
												</table>
												<!--<![endif]-->
												<!--[if gte mso 9]><div style="display:none;"><table cellspacing="0" cellpadding="0" class="EQ-00" width="550">
	<tr><td width="20" class="layout-block-padding-left">&nbsp; </td><td width="520" class="layout-block-content-cell" valign="top" align="left"><v:rect style="width:520px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg" color="#ffffff"></v:fill><v:textbox style="mso-fit-shape-to-text:true" inset="0,0,0,0"><div><div style="font-size:0"><table cellspacing="0" cellpadding="0" class="EQ-07">
	<tr><td style="font-size:1px;" width="10">&nbsp; </td><td align="left" class="EQ-05" width="500"><div class="text" style="font-size: 16px; font-family: sans-serif; line-height: 120%;"><font face="Times New Roman, sans-serif" style="line-height: 120%;">Hi ` + invite.FirstName + `,</font><div style="line-height: 120%;"><b style="font-family: sans-serif; font-size: 16px;"><br></b></div><div style="line-height: 120%;"><font face="Times New Roman, sans-serif"><b style="font-size: 16px;">` + name + `</b><span style="font-size: 16px;"> invited you to the </span><b style="font-size: 16px;">` + accountName + `</b><span style="font-size: 16px;"> account.</span></font><span style="font-family: 'Times New Roman', sans-serif;"> </span></div><div style="line-height: 120%;"><span style="font-family: 'Times New Roman', sans-serif;"><br></span></div><div style="line-height: 120%;"><span style="font-family: 'Times New Roman', sans-serif; font-size: 15px; font-style: italic;">` + invite.Message + `</span></div><div style="line-height: 120%;"><br></div></div></td><td style="font-size:1px;" width="10">&nbsp; </td>
	</tr></table></div></div></v:textbox></v:rect></td><td width="10" class="layout-block-padding-right">&nbsp; </td>
	</tr></table></div><![endif]-->
											</div>

											<div width="100%">
												<!--[if !mso 15]><!-->
												<table width="550" cellspacing="0" cellpadding="0" class="EQ-00" style="mso-hide:all;">
													<tbody>
														<tr>
															<td width="12" class="EQ-18">&nbsp;</td>
															<td background="https://cdn.skyclerk.com/emails/invite-new-user/box-bg-5.jpg" class="EQ-20" bgcolor="#ffffff" width="511">
																<table cellspacing="0" cellpadding="0" align="left" style="padding-left:10px; padding-right:10px;" class="EQ-06" width="511">
																	<tbody>
																		<tr>
																			<td valign="top" class="EQ-06" align="center" width="491">
																				<div class="layout-block-image">
																					<a href="` + url + `" target="_blank"><img width="491" height="86" alt="Sign up for Skyclerk" src="https://cdn.skyclerk.com/emails/invite-new-user/image-2.png" border="0" style="color: rgb(111, 163, 61); display: block; width: 491px; height: 86px;"
																						 class="EQ-28"></a>
																				</div>
																			</td>
																		</tr>
																	</tbody>
																</table>
															</td>
															<td width="27" class="EQ-19">&nbsp;</td>
														</tr>
													</tbody>
												</table>
												<!--<![endif]-->
												<!--[if gte mso 9]><div style="display:none;"><table cellspacing="0" cellpadding="0" class="EQ-00" width="550">
	<tr><td width="12" class="layout-block-padding-left">&nbsp; </td><td class="layout-block-content-cell" width="511"><v:rect style="width:511px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/emails/invite-new-user/box-bg-5.jpg" color="#ffffff"></v:fill><v:textbox style="mso-fit-shape-to-text:true" inset="0,0,0,0"><div><div style="font-size:0"><table cellspacing="0" cellpadding="0" align="left" style="padding-left:10px; padding-right:10px;" class="EQ-06">
	<tr><td valign="top" class="EQ-06" align="center"><div class="layout-block-image"><a href="` + url + `"><img height="86" alt="Sign up for Skyclerk" src="https://cdn.skyclerk.com/emails/invite-new-user/image-2.png" border="0" style="color: #6FA33D; display: block;"  class="EQ-28" width="491"></img><div style="width:0px;height:0px;max-height:0;max-width:0;overflow:hidden;display:none;visibility:hidden;mso-hide:all;"></div></a></div></td>
	</tr></table></div></div></v:textbox></v:rect></td><td width="27" class="layout-block-padding-right">&nbsp; </td>
	</tr></table></div><![endif]-->
											</div>

											<div width="100%">
												<!--[if !mso 15]><!-->
												<table width="550" cellspacing="0" cellpadding="0" class="EQ-00" style="mso-hide:all;">
													<tbody>
														<tr>
															<td width="20" class="EQ-22">&nbsp;</td>
															<td width="520" class="EQ-24">
																<table cellspacing="0" cellpadding="0" align="left" class="EQ-05">
																	<tbody>
																		<tr>
																			<td width="500" valign="top" align="left" background="https://cdn.skyclerk.com/emails/invite-new-user/box-bg-5.jpg" style="padding-left: 10px; padding-right: 10px;" bgcolor="#ffffff">
																				<table cellspacing="0" cellpadding="0" class="EQ-07">
																					<tbody>
																						<tr>
																							<td align="left" class="EQ-05" width="500">
																								<div class="text" style="font-size: 16px; font-family: 'Lucida Grande';">
																									<div style="line-height: 1.2;"><span style="font-family: Roboto, helvetica, &quot;Times New Roman&quot;, sans-serif; font-size: 15px;">If you’re new to Skyclerk,&nbsp;</span><a href="https://www.youtube.com/watch?v=ahQgNssTNrs"
																										 style="font-family: Roboto, helvetica, &quot;Times New Roman&quot;, sans-serif; color: rgb(111, 163, 61); font-size: 15px;" target="_blank">check out this 60 second video</a><span style="font-family: Roboto, helvetica, &quot;Times New Roman&quot;, sans-serif; font-size: 15px;">&nbsp;to
																											learn more or visit our website at <a href="https://skyclerk.com" style="color: rgb(111, 163, 61);" target="_blank">https://skyclerk.com</a>.</span></div>
																									<div style="line-height: 1.2;"><span style="font-family: &quot;Helvetica Neue&quot;, helvetica, arial, sans-serif; font-size: 15px;"><br></span></div>
																									<div style="line-height: 1.2;">
																										<font face="Roboto, Times New Roman, sans-serif">Thanks!</font>
																									</div>
																									<div style="line-height: 1.2;">
																										<font face="Roboto, Times New Roman, sans-serif">- The Skyclerk Team</font>
																									</div>
																									<div>
																										<font face="Arial"><br></font>
																									</div>
																								</div>
																							</td>
																						</tr>
																					</tbody>
																				</table>
																			</td>
																		</tr>
																	</tbody>
																</table>
															</td>
															<td width="10" class="EQ-23">&nbsp;</td>
														</tr>
													</tbody>
												</table>
												<!--<![endif]-->
												<!--[if gte mso 9]><div style="display:none;"><table cellspacing="0" cellpadding="0" class="EQ-00" width="550">
	<tr><td width="20" class="layout-block-padding-left">&nbsp; </td><td width="520" class="layout-block-content-cell" valign="top" align="left"><v:rect style="width:520px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/emails/invite-new-user/box-bg-5.jpg" color="#ffffff"></v:fill><v:textbox style="mso-fit-shape-to-text:true" inset="0,0,0,0"><div><div style="font-size:0"><table cellspacing="0" cellpadding="0" class="EQ-07">
	<tr><td style="font-size:1px;" width="10">&nbsp; </td><td align="left" class="EQ-05" width="500"><div class="text" style="font-size: 16px; font-family: sans-serif;"><div style="line-height: 120%;"><span style="font-family: 'Times New Roman', sans-serif; font-size: 15px;">If you’re new to Skyclerk, </span><a href="https://www.youtube.com/watch?v=ahQgNssTNrs" style="font-family: 'Times New Roman', sans-serif; color: #6FA33D; font-size: 15px;">check out this 60 second video</a><span style="font-family: 'Times New Roman', sans-serif; font-size: 15px;"> to learn more or visit our website at <a href="https://skyclerk.com" style="color: #6FA33D;">https://skyclerk.com</a>.</span></div><div style="line-height: 120%;"><span style="font-family: sans-serif; font-size: 15px;"><br></span></div><div style="line-height: 120%;"><font face="Times New Roman, sans-serif">Thanks!</font></div><div style="line-height: 120%;"><font face="Times New Roman, sans-serif">- The Skyclerk Team</font></div><div><font face="Arial"><br></font></div></div></td><td style="font-size:1px;" width="10">&nbsp; </td>
	</tr></table></div></div></v:textbox></v:rect></td><td width="10" class="layout-block-padding-right">&nbsp; </td>
	</tr></table></div><![endif]-->
											</div>

											<div width="100%">
												<!--[if !mso 15]><!-->
												<table width="550" cellspacing="0" cellpadding="0" class="EQ-00" style="mso-hide:all;">
													<tbody>
														<tr>
															<td width="10" class="EQ-23" style="font-size:1px;">&nbsp;</td>
															<td height="20" background="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg" class="EQ-34" bgcolor="#ffffff" width="530">
																<div class="spacer"></div>
															</td>
															<td width="10" class="EQ-23" style="font-size:1px;">&nbsp;</td>
														</tr>
													</tbody>
												</table>
												<!--<![endif]-->
												<!--[if gte mso 9]><div style="display:none;"><table cellspacing="0" cellpadding="0" class="EQ-00" width="550">
	<tr><td width="10" class="layout-block-padding-left" style="font-size:0">&nbsp; </td><td class="layout-block-content-cell" style="font-size:0;" width="530" height="20"><v:rect style="width:530px;height:20px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg" color="#ffffff"></v:fill><div class="spacer"></div></v:rect></td><td width="10" class="layout-block-padding-right" style="font-size:0">&nbsp; </td>
	</tr></table></div><![endif]-->
											</div>

											<div width="100%">
												<!--[if !mso 15]><!-->
												<table width="550" cellspacing="0" cellpadding="0" class="EQ-00" style="mso-hide:all;">
													<tbody>
														<tr>
															<td width="10" class="EQ-23">&nbsp;</td>
															<td width="530" valign="top" align="left" class="EQ-34">
																<table cellspacing="0" cellpadding="0" align="left" class="EQ-05">
																	<tbody>
																		<tr>
																			<td width="510" valign="top" align="left" background="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg" style="padding-left: 10px; padding-right: 10px;" bgcolor="#ffffff">
																				<table cellspacing="0" cellpadding="0" class="EQ-07">
																					<tbody>
																						<tr>
																							<td align="left" class="EQ-05" width="510">
																								<div class="heading" style="font-size: 16px; font-family: 'Lucida Grande';">
																									<div style="text-align: center;">
																										<font color="#bfbfbf" face="Roboto, Arial, sans-serif" style="font-size: 13px;">Skyclerk, 901 Brutscher St, D112, Newberg, OR, 97132</font>
																									</div>
																									<div style="text-align: center;">
																										<font color="#bfbfbf" face="Roboto, Arial, sans-serif" style="font-size: 13px;"><span style="caret-color: rgb(191, 191, 191);">For additional help please visit <font color="#bfbfbf" face="Roboto, Arial, sans-serif" style="font-size: 13px;">https://</font><a
																												 href="http://skyclerk.com/support" style="color: rgb(190, 191, 191);" target="_blank">skyclerk.com/support</a>.</span></font>
																									</div>
																								</div>
																							</td>
																						</tr>
																					</tbody>
																				</table>
																			</td>
																		</tr>
																	</tbody>
																</table>
															</td>
															<td width="10" class="EQ-23">&nbsp;</td>
														</tr>
													</tbody>
												</table>
												<!--<![endif]-->
												<!--[if gte mso 9]><div style="display:none;"><table cellspacing="0" cellpadding="0" class="EQ-00" width="550">
	<tr><td width="10" class="layout-block-padding-left">&nbsp; </td><td width="530" valign="top" align="left" class="layout-block-content-cell"><v:rect style="width:530px;" stroke="f"><v:fill type="tile" src="https://cdn.skyclerk.com/emails/invite-new-user/box-bg.jpg" color="#ffffff"></v:fill><v:textbox style="mso-fit-shape-to-text:true" inset="0,0,0,0"><div><div style="font-size:0"><table cellspacing="0" cellpadding="0" align="left" class="EQ-05">
	<tr><td style="font-size:1px;" width="10">&nbsp; </td><td width="510" valign="top" align="left"><div class="heading" style="font-size: 16px; font-family: sans-serif;"><div style="text-align: center;"><font color="#bfbfbf" face="Arial, sans-serif" style="font-size: 13px;">Skyclerk, 901 Brutscher St, D112, Newberg, OR, 97132</font></div><div style="text-align: center;"><font color="#bfbfbf" face="Arial, sans-serif" style="font-size: 13px;"><span style="caret-color: #BFBFBF;">For additional help please visit <font color="#bfbfbf" face="Arial, sans-serif" style="font-size: 13px;">https://</font><a href="http://skyclerk.com/support" style="color: #BEBFBF;">skyclerk.com/support</a>.</span></font></div></div></td><td style="font-size:1px;" width="10">&nbsp; </td>
	</tr></table></div></div></v:textbox></v:rect></td><td width="10" class="layout-block-padding-right">&nbsp; </td>
	</tr></table></div><![endif]-->
											</div>

										</div>
										<!--End of layout container -->

									</td>
									<td width="30" class="EQ-01">&nbsp;
										<!--Right page bg show-thru -->
									</td>
								</tr>
							</tbody>
						</table>
						<!--email-body -->

					</td>
				</tr>
			</tbody>
		</table>
		<!--page-wrap -->


	</body>

	</html>


	`
}

/* End File */
