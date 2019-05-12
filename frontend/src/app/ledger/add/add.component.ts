//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input } from '@angular/core';

@Component({
	selector: 'app-ledger-add',
	templateUrl: './add.component.html'
})
export class AddComponent implements OnInit {
	@Input() type: string = "income";

	showAddContact: boolean = false;

	//
	// Constructor
	//
	constructor() { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Add contact click
	//
	addContactToggle() {
		if (this.showAddContact) {
			this.showAddContact = false;
		} else {
			this.showAddContact = true;
		}
	}

}

/* End File */
