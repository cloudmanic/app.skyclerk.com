//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment-timezone';
import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { ReportService } from 'src/app/services/report.service';

declare var Pikaday: any;

@Component({
	selector: 'app-dashboard-graphs',
	templateUrl: './graphs.component.html',
})

export class GraphsComponent implements OnInit {
	showFilter: boolean = false;

	// Setup date pickers
	endDate: Date = new Date();
	startDate: Date = moment(moment().format('YYYY-01-01')).toDate();
	@ViewChild('endDateField') endDateField: ElementRef;
	@ViewChild('endDateTrigger') endDateTrigger: ElementRef;
	@ViewChild('startDateField') startDateField: ElementRef;
	@ViewChild('startDateTrigger') startDateTrigger: ElementRef;

	//
	// Constructor
	//
	constructor(public reportService: ReportService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Setup date pickers.
		this.setupDatePickers();

		// Build the page
		this.refreshPageData();
	}

	//
	// Load the page data.
	//
	refreshPageData() {
		console.log(this.startDate);
		console.log(this.endDate);
	}

	//
	// Setup date pickers
	//
	setupDatePickers() {
		// Setup start date picker.
		new Pikaday({
			defaultDate: this.startDate,
			field: this.startDateField.nativeElement,
			trigger: this.startDateTrigger.nativeElement,
			onSelect: (date: Date) => {
				this.startDate = date;
				this.refreshPageData();
			}
		});

		// Setup end date picker.
		new Pikaday({
			defaultDate: this.endDate,
			field: this.endDateField.nativeElement,
			trigger: this.endDateTrigger.nativeElement,
			onSelect: (date: Date) => {
				this.endDate = date;
				this.refreshPageData();
			}
		});
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
}

/* End File */
