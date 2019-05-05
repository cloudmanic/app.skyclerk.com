//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { AuthService } from 'src/app/services/auth.service';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { MeService } from 'src/app/services/me.service';

@Component({
	selector: 'app-auth-login',
	templateUrl: './login.component.html'
})

export class LoginComponent implements OnInit {
	errMsg: string = "";
	email: string = "";
	password: string = "";

	//
	// Construct.
	//
	constructor(public authService: AuthService, public meService: MeService, public router: Router) { }

	//
	// OnInit...
	//
	ngOnInit() { }

	//
	// Submit login request
	//
	doLogin() {
		// Clear old error
		this.errMsg = "";

		// Make the the HTTP request:
		this.authService.login(this.email, this.password).subscribe(
			// Success - Redirect to dashboard.
			() => {
				// Get the user so we can set the default account id
				this.meService.get().subscribe(res => {
					localStorage.setItem('account_id', res.Accounts[0].Id.toString());
					this.router.navigate(['/']);
				});
			},

			// Error
			(err: HttpErrorResponse) => {
				this.errMsg = err.error.error;
			}
		);
	}

}

/* End File */
