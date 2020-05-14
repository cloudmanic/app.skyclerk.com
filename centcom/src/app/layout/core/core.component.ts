//
// Date: 2019-11-04
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
	selector: 'app-core',
	templateUrl: './core.component.html'
})

export class CoreComponent {
	tab: string = "";
	pageHeading: string = "";
	profileShow: boolean = false;

	//
	// Constructor
	//
	constructor(public http: HttpClient, public route: ActivatedRoute, public router: Router) {
		// Figure out which tab this is.
		this.router.events.subscribe((val) => {
			this.tab = this.route.snapshot.firstChild.data.tab;
			this.pageHeading = this.route.snapshot.firstChild.data.pageHeading;
		});

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
