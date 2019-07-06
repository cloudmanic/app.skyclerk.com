//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { ReportService } from 'src/app/services/report.service';

@Component({
	selector: 'app-dashboard-graphs',
	templateUrl: './graphs.component.html',
})

export class GraphsComponent implements OnInit {
	showFilter: boolean = false;
	editStartDate: boolean = false;

	//
	// Constructor
	//
	constructor(public reportService: ReportService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// // Load page data
		// this.loadActivity();
		//
		// // Build charts
		// this.buildChart1();
		// this.buildChart2();
	}

	//
	// Toggle filter
	//
	toggleFilter() {
		if (this.showFilter) {
			this.showFilter = false;
		} else {
			this.showFilter = true;
		}
	}

	//
	// Toggle start date
	//
	toggleEditStartDate() {
		if (this.editStartDate) {
			this.editStartDate = false;
		} else {
			this.editStartDate = true;
		}
	}
}

/* End File */
