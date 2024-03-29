//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { map } from "rxjs/operators";
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';

@Injectable({
	providedIn: 'root'
})

export class AuthService {
	//
	// Constructor.
	//
	constructor(private http: HttpClient) { }

	//
	// Log user in.
	//
	login(email: string, pass: string): Observable<LoginResponse> {
		// Build POST
		let post = {
			username: email,
			password: pass,
			client_id: environment.client_id,
			grant_type: 'password'
		}

		// Send Request to BE.
		return this.http.post<LoginResponse>(environment.app_server + '/oauth/token', post)
			.pipe(map(res => {
				// Store access token in local storage.
				localStorage.setItem('user_id', res["user_id"].toString());
				localStorage.setItem('access_token', res["access_token"]);
				localStorage.setItem('user_email', email);

				return {
					user_id: res["user_id"],
					access_token: res["access_token"]
				};
			}));
	}

	//
	// Forgot password request
	//
	forgotPassword(email: string): Observable<boolean> {
		// Build POST
		let post = {
			email: email
		}

		// Send Request to BE.
		return this.http.post<boolean>(environment.app_server + '/forgot-password', post)
			.pipe(map(() => {
				return true;
			}));
	}

	//
	// Reset password request
	//
	resetPassword(password: string, hash: string): Observable<boolean> {
		// Build POST
		let post = {
			password: password,
			hash: hash
		}

		// Send Request to BE.
		return this.http.post<boolean>(environment.app_server + '/reset-password', post)
			.pipe(map(() => {
				return true;
			}));
	}
}

// Response from a register request.
export interface RegisterResponse {
	user_id: number,
	access_token: string
}

// Response from a login request
export interface LoginResponse {
	user_id: number,
	access_token: string
}

/* End File */
