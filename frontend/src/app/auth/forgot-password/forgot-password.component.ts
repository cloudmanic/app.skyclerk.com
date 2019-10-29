//
// Date: 2019-10-28
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AuthService } from 'src/app/services/auth.service';
import { MeService } from 'src/app/services/me.service';
import { Router } from '@angular/router';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';
import { HttpErrorResponse } from '@angular/common/http';

const pageTitle: string = environment.title_prefix + "Forgot Password";

@Component({
	selector: 'app-forgot-password',
	templateUrl: './forgot-password.component.html',
})

export class ForgotPasswordComponent implements OnInit {
	errMsg: string = "";
	email: string = "";
	successMsg: string = "";

	//
	// Construct.
	//
	constructor(public authService: AuthService, public meService: MeService, public router: Router, public titleService: Title) { }

	//
	// OnInit...
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);
	}

	//
	// Submit forgot password request.
	//
	submit() {
		// Clear old error
		this.errMsg = "";
		this.successMsg = "";

		// Make the the HTTP request:
		this.authService.forgotPassword(this.email).subscribe(
			// Success - Redirect to dashboard.
			() => {
				this.successMsg = "Please check your email for further instructions.";
				this.email = "";
			},

			// Error
			(err: HttpErrorResponse) => {
				this.errMsg = err.error.error;
			}
		);
	}

}

/* End File */
