//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';

@Component({
	selector: 'app-settings-account-company-name',
	templateUrl: './company-name.component.html'
})

export class CompanyNameComponent {
	editMode: boolean = false;

	//
	// Constructor
	//
	constructor() { }

	//
	// Toggle to edit mode.
	//
	editModeToggle() {
		this.editMode = !this.editMode;
	}
}

/* End File */
