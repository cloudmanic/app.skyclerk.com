//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Router } from '@angular/router';
import { MeService } from './me.service';

@Injectable({
	providedIn: 'root'
})

export class PingService {
	//
	// Constructor.
	//
	constructor(private http: HttpClient, private router: Router, private meService: MeService) { }

	//
	// Start ping to BE Server.
	//
	startPing() {
		// Ping server at start too
		this.PingServer();

		// Ajax call to ping server every 10 seconds
		setInterval(() => { this.PingServer(); }, 10000);
	}

	// Ping to make sure our access token is still good.
	// If not redirect back to login. Also the server
	// uses this as an opportunity to collect some stats.
	//
	PingServer() {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/ping`;

		if (!accountId) {
			return;
		}

		// Send ping to server
		this.http.get(url).subscribe(
			// Success
			data => {

				// We do not redirect from - /paywall
				if (window.location.pathname == "/paywall") {
					return;
				}

				// We do not redirect from - /paywall/plans
				if (window.location.pathname == "/paywall/plans") {
					return;
				}

				// We do not redirect from - /paywall/payment
				if (window.location.pathname == "/paywall/payment") {
					return;
				}

				// We do not redirect from - /paywall/success
				if (window.location.pathname == "/paywall/success") {
					return;
				}


				// // We do not redirect from - /register
				// if (window.location.pathname == "/register") {
				// 	return;
				// }

				// We do not redirect from - /forgot-password
				if (window.location.pathname == "/login") {
					return;
				}

				// We do not redirect from - /forgot-password
				if (window.location.pathname == "/forgot-password") {
					return;
				}

				// We do not redirect from - /reset-password
				if (window.location.pathname == "/reset-password") {
					return;
				}

				// Delinquent status - Means the person is not current on payment.
				if (data["status"] == "delinquent") {
					this.router.navigate(['/paywall']);
					return;
				}

				// Expired status - Means the person's free trial has expired.
				if (data["status"] == "expired") {
					this.router.navigate(['/paywall']);
					return;
				}

				// Logout status
				if (data["status"] == "logout") {
					this.meService.logout();
					this.router.navigate(['/login']);
					return;
				}

			},

			// Error
			(err: HttpErrorResponse) => {
				// We do not redirect from - /paywall
				if (window.location.pathname == "/paywall") {
					return;
				}

				// We do not redirect from - /paywall/plans
				if (window.location.pathname == "/paywall/plans") {
					return;
				}

				// We do not redirect from - /paywall/payment
				if (window.location.pathname == "/paywall/payment") {
					return;
				}

				// We do not redirect from - /paywall/success
				if (window.location.pathname == "/paywall/success") {
					return;
				}


				// // We do not redirect from - /register
				// if (window.location.pathname == "/register") {
				// 	return;
				// }

				// We do not redirect from - /forgot-password
				if (window.location.pathname == "/login") {
					return;
				}

				// We do not redirect from - /forgot-password
				if (window.location.pathname == "/forgot-password") {
					return;
				}

				// We do not redirect from - /reset-password
				if (window.location.pathname == "/reset-password") {
					return;
				}

				// Error....
				if (err.error instanceof Error) {
					// A client-side or network error occurred. Handle it accordingly.
					console.log('An error occurred:', err.error);
				} else {
					// Logged out user.
					if (err.status == 401) {
						this.meService.logout();
						this.router.navigate(['/login']);
						return;
					}

					// // Access token mostly not good.
					// // If the error is blank it often means the
					// // server is down.
					// this.router.navigate(['/logout']);
					// if (err.error.error && (err.error.error.length > 0)) {
					// 	this.router.navigate(['/logout']);
					// }
				}

			}

		);

	}
}

/* End File */
