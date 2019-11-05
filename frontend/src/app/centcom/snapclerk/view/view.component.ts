//
// Date: 2019-11-04
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//


import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';

@Component({
	selector: 'app-view',
	templateUrl: './view.component.html'
})

export class ViewComponent implements OnInit {
	contact: string = "";
	contactsResults: string[] = ["Spicer was here", "I was here", "jane was a runner"];

	//
	// Constructor
	//
	constructor(public http: HttpClient) { }

	//
	// ngOnInit
	//
	ngOnInit() {


	}

	//
	// onContactChange do search
	//
	onContactChange() {
		console.log(this.contact);
	}

	//
	// onContactSelect
	//
	onContactSelect(result: string) {
		this.contact = result;
		this.contactsResults = [];
	}
}

/* End File */
