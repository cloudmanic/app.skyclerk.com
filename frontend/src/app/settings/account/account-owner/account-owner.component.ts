//
// Date: 2019-08-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component } from '@angular/core';

@Component({
	selector: 'app-settings-account-account-owner',
	templateUrl: './account-owner.component.html'
})

export class AccountOwnerComponent {
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
