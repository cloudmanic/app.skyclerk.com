//
// Date: 2019-05-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input } from '@angular/core';
import { BaseComponent } from '../base/base.component';
import { ContactService } from 'src/app/services/contact.service';

@Component({
	selector: 'app-contacts-add-edit-pop-out',
	templateUrl: './add-edit-pop-out.component.html'
})

export class AddEditPopOutComponent extends BaseComponent implements OnInit {
	@Input() show: boolean = false;

	showMore: boolean = false;

	//
	// Constructor
	//
	constructor(public contactService: ContactService) {
		super(contactService);
	}

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Show more
	//
	showMoreClick() {
		this.showMore = true;
	}

	//
	// Cancel popup.
	//
	cancelPopUp() {
		this.show = false;
		this.showMore = false;
	}
}

/* End File */
