//
// Date: 2019-04-27
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { map } from "rxjs/operators";
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';
import { TrackService } from './track.service';
import { User } from '../models/user.model';

@Injectable({
	providedIn: 'root'
})

export class UserService {
	//
	// Constructor
	//
	constructor(private http: HttpClient, private trackService: TrackService) { }

	//
	// Get users
	//
	get(): Observable<User[]> {
		let accountId = localStorage.getItem('account_id');
		let url = environment.app_server + '/api/v3/' + accountId + '/users';
		return this.http.get<User[]>(url)
			.pipe(map(res => res.map(res => new User().deserialize(res))));
	}

	//
	// Returns a list of non-expired invites.
	//
	getInvites(): Observable<User[]> {
		let accountId = localStorage.getItem('account_id');

		return this.http.get<User[]>(`${environment.app_server}/api/v3/${accountId}/users/invite`)
			.pipe(map(res => res.map(res => new User().deserialize(res))));
	}

	//
	// Invite a new user to the application
	//
	invite(first: string, last: string, email: string, message: string): Observable<boolean> {
		let accountId = localStorage.getItem('account_id');

		return this.http.post<boolean>(`${environment.app_server}/api/v3/${accountId}/users/invite`, { first_name: first, last_name: last, email: email, message: message })
			.pipe(map(() => {
				// Track event.
				this.trackService.event('user-invite', { app: "web", "accountId": accountId });

				return true;
			}));
	}

	//
	// Delete a Invite
	//
	deleteInvite(id: number): Observable<Boolean> {
		let accountId = localStorage.getItem('account_id');

		return this.http.delete<Boolean>(`${environment.app_server}/api/v3/${accountId}/user-invite/${id}`, {})
			.pipe(map(() => {
				// Track event.
				this.trackService.event('user-invite-delete', { app: "web", "accountId": accountId });

				return true;
			}));
	}

	//
	// Delete a User
	//
	deleteUser(id: number): Observable<Boolean> {
		let accountId = localStorage.getItem('account_id');

		return this.http.delete<Boolean>(`${environment.app_server}/api/v3/${accountId}/users/${id}`, {})
			.pipe(map(() => {
				// Track event.
				this.trackService.event('user-delete', { app: "web", "accountId": accountId });

				return true;
			}));
	}

	//
	// //
	// // Update a label
	// //
	// update(label: Label): Observable<Label> {
	// 	let accountId = localStorage.getItem('account_id');
	//
	// 	let put = {
	// 		name: label.Name,
	// 		account_id: Number(accountId)
	// 	}
	//
	// 	return this.http.put<Label>(`${environment.app_server}/api/v3/${accountId}/labels/${label.Id}`, put)
	// 		.pipe(map(res => {
	// 			let lb = new Label().deserialize(res);
	//
	// 			// Track event.
	// 			this.trackService.event('label-update', { app: "web", "accountId": accountId });
	//
	// 			return lb;
	// 		}));
	// }
}


/* End File */
