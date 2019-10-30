//
// Date: 2019-10-28
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { Params, Router, ActivatedRoute } from '@angular/router';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';

const pageTitle: string = environment.title_prefix + "Reset Password";

@Component({
	selector: 'app-reset-password',
	templateUrl: './reset-password.component.html'
})

export class ResetPasswordComponent implements OnInit {
	hash: string = "";
	errorMsg: string = "";
	password: string = "";
	passwordConfirm: string = "";
	successMsg: string = "";

	constructor(private router: Router, private activatedRoute: ActivatedRoute, private titleService: Title) { }

	//
	// OnInit...
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// subscribe to router event
		this.activatedRoute.queryParams.subscribe((params: Params) => {
			this.hash = params['hash'];
		});

	}

	//
	// Reset submit.
	//
	submit() {
		// Make sure we have a password.
		if (this.password.length == 0) {
			this.errorMsg = "Opps, please submit a password.";
			return;
		}

		// Make sure the passwords match.
		if (this.password != this.passwordConfirm) {
			this.errorMsg = "Opps, the passwords do not match.";
			return;
		}

		// Clear old error
		this.errorMsg = "";
		this.successMsg = "";

		// Make the the HTTP request:
		this.authService.resetPassword(this.password, this.hash).subscribe(
			// Success - Redirect to dashboard.
			() => {
				this.successMsg = "Please check your email for further instructions.";
				this.password = "";
				this.passwordConfirm = "";
			},

			// Error
			(err: HttpErrorResponse) => {
				this.errorMsg = err.error.error;
			}
		);

		// // First make sure the passwords match.
		// if (form.value.password != form.value.password_again) {
		// 	this.errorMsg = "Opps, the passwords do not match.";
		// 	return;
		// }
		//
		// // Clear post error.
		// this.errorMsg = "";
		//
		// // Update submit button
		// this.submitBtn = "Saving...";
		//
		// // Add the hash to the post.
		// form.value.hash = this.hash;
		//
		// // Make the the HTTP request:
		// this.http.post(environment.app_server + '/reset-password', form.value).subscribe(
		//
		// 	// Success redirect to login
		// 	data => {
		// 		this.router.navigate(['/login'], { queryParams: { success: "Your password was successfully reset." } });
		// 	},
		//
		// 	// Error
		// 	(err: HttpErrorResponse) => {
		//
		// 		// Change button back.
		// 		this.submitBtn = "Reset Password";
		//
		// 		if (err.error instanceof Error) {
		// 			console.log('A client-side error occurred:', err.error.message);
		// 		} else {
		// 			this.errorMsg = err.error.error;
		// 		}
		//
		// 	}
		//
		// );

	}


}

/* End File */
