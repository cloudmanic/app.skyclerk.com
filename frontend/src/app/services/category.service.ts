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
import { Category } from '../models/category.model';
import { TrackService } from './track.service';

@Injectable({
	providedIn: 'root'
})

export class CategoryService {
	//
	// Constructor
	//
	constructor(private http: HttpClient, private trackService: TrackService) { }

	//
	// Get categories
	//
	get(type: string): Observable<Category[]> {
		let accountId = localStorage.getItem('account_id');
		let url = `${environment.app_server}/api/v3/${accountId}/categories?type=${type}`;
		return this.http.get<Category[]>(url)
			.pipe(map(res => res.map(res => new Category().deserialize(res))));
	}

	//
	// Create a new category
	//
	create(category: Category): Observable<Category> {
		let accountId = localStorage.getItem('account_id');
		category.AccountId = Number(accountId);

		return this.http.post<Category>(`${environment.app_server}/api/v3/${accountId}/categories`, new Category().serialize(category))
			.pipe(map(res => {
				let con = new Category().deserialize(res);

				// Track event.
				this.trackService.event('category-create', { app: "web", "accountId": accountId });

				return con;
			}));
	}

	//
	// Update a category
	//
	update(category: Category): Observable<Category> {
		let accountId = localStorage.getItem('account_id');

		// Set type. 1 = expense, 2 = income
		let type = "1";
		if (category.Type == "income") {
			type = "2"
		}

		let put = {
			type: type,
			name: category.Name,
			account_id: Number(accountId)
		}

		return this.http.put<Category>(`${environment.app_server}/api/v3/${accountId}/categories/${category.Id}`, put)
			.pipe(map(res => {
				let cat = new Category().deserialize(res);

				// Track event.
				this.trackService.event('category-update', { app: "web", "accountId": accountId });

				return cat;
			}));
	}

	//
	// Delete a category
	//
	delete(category: Category): Observable<Boolean> {
		let accountId = localStorage.getItem('account_id');
		category.AccountId = Number(accountId);

		return this.http.delete<Boolean>(`${environment.app_server}/api/v3/${accountId}/categories/${category.Id}`, {})
			.pipe(map(() => {
				// Track event.
				this.trackService.event('category-delete', { app: "web", "accountId": accountId });

				return true;
			}));
	}
}


/* End File */
