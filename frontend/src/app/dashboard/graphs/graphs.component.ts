//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import * as Highcharts from 'highcharts';
import * as moment from 'moment-timezone';
import { AngularCsv } from 'angular7-csv/dist/Angular-csv'
import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { ReportService, PnlNameAmount } from 'src/app/services/report.service';
import { Subject } from 'rxjs';
import { MeService } from 'src/app/services/me.service';

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
	destory: Subject<boolean> = new Subject<boolean>();

	// Setup date pickers
	endDate: Date = new Date();
	startDate: Date = moment(moment().format('YYYY-01-01')).toDate();
	@ViewChild('endDateField') endDateField: ElementRef;
	@ViewChild('endDateTrigger') endDateTrigger: ElementRef;
	@ViewChild('startDateField') startDateField: ElementRef;
	@ViewChild('startDateTrigger') startDateTrigger: ElementRef;

	// Setup chart options
	chartOptions: any = {
		chart: { type: 'column' },

		title: { text: '' },

		credits: { enabled: false },

		rangeSelector: { enabled: false },

		scrollbar: { enabled: false },

		navigator: { enabled: false },

		legend: { enabled: false },

		tooltip: {
			formatter: function() {
				// TODO(spicer): Manage different currencies
				return this.x + ': $' + Highcharts.numberFormat(this.y, 0, '.', ',');
			}
		},

		yAxis: {
			title: { text: '' },

			labels: {
				formatter: function() {
					// TODO(spicer): Manage different currencies
					return '$' + Highcharts.numberFormat(this.value, 0, '.', ',');
				}
			}
		},

		xAxis: {
			categories: []
		},

		series: [
			{ name: "", color: "#757575", data: [] }
		]
	}

	//
	// Constructor
	//
	constructor(public reportService: ReportService, public meService: MeService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Setup date pickers.
		this.setupDatePickers();

		// Build the page
		this.refreshPageData();

		// Listen for account changes.
		this.meService.accountChange.takeUntil(this.destory).subscribe(() => {
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

			// Build Graph
			let cats = [];
			let data = [];

			for (let i = 0; i < res.length; i++) {
				// Set X-Axis
				cats.push(res[i].Name);

				// Set color
				let color = "#537b37";

				// Set Y-Axis
				data.push({ color: color, y: res[i].Amount });
			}

			// Rebuilt the chart
			this.chartOptions.series[0].name = this.nameTitle;
			this.chartOptions.series[0].data = data;
			this.chartOptions.xAxis.categories = cats;
			Highcharts.chart('chart', this.chartOptions);
		});
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

			// Build Graph
			let cats = [];
			let data = [];

			for (let i = 0; i < res.length; i++) {
				// Set X-Axis
				cats.push(res[i].Name);

				// Set color
				let color = "#bb4626";

				// Set Y-Axis
				data.push({ color: color, y: (res[i].Amount * -1) });
			}

			// Rebuilt the chart
			this.chartOptions.series[0].name = this.nameTitle;
			this.chartOptions.series[0].data = data;
			this.chartOptions.xAxis.categories = cats;
			Highcharts.chart('chart', this.chartOptions);
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

			// Build Graph
			let cats = [];
			let data = [];

			for (let i = 0; i < res.length; i++) {
				// Set X-Axis
				cats.push(res[i].Name);

				// Set color
				let color = "#537b37";

				if (res[i].Amount < 0) {
					color = "#bb4626";
				}

				// Set Y-Axis
				data.push({ color: color, y: res[i].Amount });
			}

			// Rebuilt the chart
			this.chartOptions.series[0].name = this.nameTitle;
			this.chartOptions.series[0].data = data;
			this.chartOptions.xAxis.categories = cats;
			Highcharts.chart('chart', this.chartOptions);
		});
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

			// Build Graph
			let cats = [];
			let data = [];

			for (let i = 0; i < res.length; i++) {
				// Set X-Axis
				cats.push(res[i].Name);

				// Set color
				let color = "#537b37";

				if (res[i].Amount < 0) {
					color = "#bb4626";
				}

				// Set Y-Axis
				data.push({ color: color, y: res[i].Amount });
			}

			// Rebuilt the chart
			this.chartOptions.series[0].name = this.nameTitle;
			this.chartOptions.series[0].data = data;
			this.chartOptions.xAxis.categories = cats;
			Highcharts.chart('chart', this.chartOptions);
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

	//
	// Export CSV
	//
	doExportCSV() {
		let options = {
			fieldSeparator: ',',
			quoteStrings: '"',
			decimalseparator: '.',
			headers: [this.nameTitle, 'Amount'],
			showTitle: false,
			useBom: true,
			removeNewLines: false,
			keys: []
		};

		// Download CSV to browser.
		new AngularCsv(this.nameAmount, 'skyclerk-' + this.type, options);
	}
}

/* End File */
