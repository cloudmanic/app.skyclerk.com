//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input } from '@angular/core';

@Component({
	selector: 'app-dashboard-sub-nav',
	templateUrl: './sub-nav.component.html'
})

export class SubNavComponent implements OnInit {
	@Input() current: string = "";

	constructor() { }

	ngOnInit() {
	}

}

/* End File */
