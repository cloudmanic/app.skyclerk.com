//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';

const pageTitle: string = environment.title_prefix + "Settings Account";

@Component({
	selector: 'app-settings-account',
	templateUrl: './account.component.html'
})
export class AccountComponent implements OnInit {

	//
	// Construct.
	//
	constructor(public titleService: Title) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);
	}

}

/* End File */
