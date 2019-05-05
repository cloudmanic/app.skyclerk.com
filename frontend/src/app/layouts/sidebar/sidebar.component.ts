//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input } from '@angular/core';

@Component({
	selector: 'app-layouts-sidebar',
	templateUrl: './sidebar.component.html'
})
export class SidebarComponent implements OnInit {
	@Input() current: string = "";

	//
	// Constructor
	//
	constructor() { }

	//
	// NgOnInit
	//
	ngOnInit() { }

}

/* End File */
