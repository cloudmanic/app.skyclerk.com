//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import * as moment from 'moment-timezone';
import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { ReportService, PnlNameAmount } from 'src/app/services/report.service';
import { Subject } from 'rxjs';
import { Title } from '@angular/platform-browser';
import { AccountService } from 'src/app/services/account.service';

declare var Pikaday: any;

@Component({
	selector: 'app-dashboard-reports',
	templateUrl: './reports.component.html'
})

export class ReportsComponent implements OnInit {
	pageTitle: string = "Skyclerk | Dashboard Reports";
	showFilter: boolean = false;
	type: string = "Income Statement";
	nameAmount: PnlNameAmount[] = [];
	destory: Subject<boolean> = new Subject<boolean>();

	// Setup date pickers
	endDate: Date = new Date();
	startDate: Date = moment(moment().format('YYYY-01-01')).toDate();
	@ViewChild('endDateField') endDateField: ElementRef;
	@ViewChild('endDateTrigger') endDateTrigger: ElementRef;
	@ViewChild('startDateField') startDateField: ElementRef;
	@ViewChild('startDateTrigger') startDateTrigger: ElementRef;

	// Income statement vars.
	incomeTotal: number = 0.00;
	expenseTotal: number = 0.00;

	//
	// Constructor
	//
	constructor(public reportService: ReportService, public accountService: AccountService, private titleService: Title) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(this.pageTitle);

		// Setup date pickers.
		this.setupDatePickers();

		// Build the page
		this.refreshPageData();

		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.refreshPageData();
		});
	}

	//
	// OnDestroy
	//
	ngOnDestroy() {
		this.destory.next();
		this.destory.complete();
	}

	//
	// Grab page data.
	//
	refreshPageData() {
		// Which report type is this?
		switch (this.type) {
			case "Income Statement":
				this.loadIncomeStatement();
				break;
		}
	}

	//
	// Load income statement.
	//
	loadIncomeStatement() {
		this.nameAmount = [];
		this.incomeTotal = 0.00;
		this.expenseTotal = 0.00;

		// AJAX call to get data.
		this.reportService.getProfitLossByCategory(this.startDate, this.endDate, "asc").subscribe(res => {
			this.nameAmount = res;

			// Figure out the income and expense totals.
			for (let i = 0; i < this.nameAmount.length; i++) {
				if (this.nameAmount[i].Amount > 0) {
					this.incomeTotal += this.nameAmount[i].Amount;
				} else {
					this.expenseTotal += (this.nameAmount[i].Amount * -1);
				}
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
	// A filter for only income entries.
	//
	incomeEntries(row: PnlNameAmount) {
		return row.Amount >= 0;
	}

	//
	// A filter for only expense entries.
	//
	expenseEntries(row: PnlNameAmount) {
		return row.Amount < 0;
	}
}

/* End File */
