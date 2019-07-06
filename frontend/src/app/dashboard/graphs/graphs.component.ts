//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment-timezone';
import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { ReportService, PnlNameAmount } from 'src/app/services/report.service';

declare var Pikaday: any;

@Component({
	selector: 'app-dashboard-graphs',
	templateUrl: './graphs.component.html',
})

export class GraphsComponent implements OnInit {
	showFilter: boolean = false;
	nameTitle: string = "Category";
	type: string = "Profit & Loss by Category";
	nameAmount: PnlNameAmount[] = [];

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
		this.nameAmount = [];

		// Which report type is this?
		switch (this.type) {
			case "Income by Contact":
				this.loadIncomeByContact();
				break;

			case "Expense by Contact":
				this.loadExpenseByContact();
				break;

			case "Profit & Loss by Label":
				this.loadProfitLossByLabel();
				break;

			case "Profit & Loss by Category":
				this.loadProfitLossByCategory();
				break;
		}
	}

	//
	// loadIncomeByContact
	//
	loadIncomeByContact() {
		// Set titles
		this.nameTitle = "Contact";

		// AJAX call to get data.
		this.reportService.getIncomeByContact(this.startDate, this.endDate, "asc").subscribe(res => {
			this.nameAmount = res;
		})
	}

	//
	// loadExpenseByContact
	//
	loadExpenseByContact() {
		// Set titles
		this.nameTitle = "Contact";

		// AJAX call to get data.
		this.reportService.getExpenseByContact(this.startDate, this.endDate, "asc").subscribe(res => {
			this.nameAmount = res;
		})
	}

	//
	// loadProfitLossByCategory
	//
	loadProfitLossByCategory() {
		// Set titles
		this.nameTitle = "Category";

		// AJAX call to get data.
		this.reportService.getProfitLossByCategory(this.startDate, this.endDate, "asc").subscribe(res => {
			this.nameAmount = res;
		})
	}

	//
	// loadProfitLossByLabel
	//
	loadProfitLossByLabel() {
		// Set titles
		this.nameTitle = "Label";

		// AJAX call to get data.
		this.reportService.getProfitLossByLabel(this.startDate, this.endDate, "asc").subscribe(res => {
			this.nameAmount = res;
		})
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
