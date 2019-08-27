//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';

const pageTitle: string = environment.title_prefix + "Settings Categories / Labels";

@Component({
	selector: 'app-settings-categories-labels',
	templateUrl: './categories-labels.component.html'
})
export class CategoriesLabelsComponent implements OnInit {

	//
	// Constructor
	//
	constructor(private titleService: Title) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);
	}

}

/* End Files */
