//
// Date: 2019-11-04
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';

@Component({
	selector: 'app-core',
	templateUrl: './core.component.html'
})

export class CoreComponent {
	profileShow: boolean = false;

	//
	// Constructor
	//
	constructor(public http: HttpClient) {
		// Make sure I am allowed to be here.
		this.ping();
	}

	//
	// Ping backend server. See if we are allowed to be here.
	//
	ping() {
		this.http.get(environment.app_server + '/api/admin/ping').subscribe(
			// Success
			() => {
				//console.log(res);
			},

			// Error
			(err) => {
				console.log(err.status);
				console.log(err.error.error);
			}
		);
	}

}

/* End File */
