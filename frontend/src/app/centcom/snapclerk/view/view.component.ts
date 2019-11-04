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

	//
	// Constructor
	//
	constructor(public http: HttpClient) { }

	//
	// ngOnInit
	//
	ngOnInit() {


	}


}

/* End File */
