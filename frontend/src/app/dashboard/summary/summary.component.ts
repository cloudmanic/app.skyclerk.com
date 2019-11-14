//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import * as moment from 'moment-timezone';
import * as Highcharts from 'highcharts';
import { Component, OnInit } from '@angular/core';
import { ActivityService } from 'src/app/services/activity.service';
import { Activity } from 'src/app/models/activity.model';
import { ReportService } from 'src/app/services/report.service';
import { Subject } from 'rxjs';
import { Title } from '@angular/platform-browser';
import { AccountService } from 'src/app/services/account.service';
import { LedgerService } from 'src/app/services/ledger.service';

@Component({
	selector: 'app-dashboard-summary',
	templateUrl: './summary.component.html'
})

export class SummaryComponent implements OnInit {
	firstRun: boolean = false;
	normalPage: boolean = false;
	pageTitle: string = "Skyclerk | Dashboard";
	activity: Activity[] = [];
	destory: Subject<boolean> = new Subject<boolean>();

	// Setup chart options - Chart 1
	chartOptions1: any = {
		chart: { type: 'column' },

		title: { text: '' },

		credits: { enabled: false },

		rangeSelector: { enabled: false },

		scrollbar: { enabled: false },

		navigator: { enabled: false },

		legend: { enabled: true },

		time: {
			getTimezoneOffset: function(timestamp: any) {
				// America/Los_Angeles
				// America/New_York
				let timezoneOffset = -moment.tz(timestamp, 'America/Los_Angeles').utcOffset();
				return timezoneOffset;
			}
		},

		tooltip: {
			formatter: function() {
				// TODO(spicer): Manage different currencies
				return Highcharts.dateFormat('%b %y', this.x) + ': $' + Highcharts.numberFormat(this.y, 0, '.', ',');
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
			type: 'datetime',
			labels: {
				formatter: function() {
					return Highcharts.dateFormat('%b %y', this.value);
				}
			}
		},

		series: [
			{ name: 'Income', color: "#537b37", data: [] },
			{ name: 'Expense', color: "#bb4626", data: [] }
		]
	}

	// Setup chart options - Chart 2
	chartOptions2: any = {
		chart: { type: 'line' },

		title: { text: '' },

		credits: { enabled: false },

		rangeSelector: { enabled: false },

		scrollbar: { enabled: false },

		navigator: { enabled: false },

		legend: { enabled: false },

		time: {
			getTimezoneOffset: function(timestamp: any) {
				// America/Los_Angeles
				// America/New_York
				let timezoneOffset = -moment.tz(timestamp, 'America/Los_Angeles').utcOffset();
				return timezoneOffset;
			}
		},

		tooltip: {
			formatter: function() {
				// TODO(spicer): Manage different currencies
				return Highcharts.dateFormat('%b %y', this.x) + ': $' + Highcharts.numberFormat(this.y, 0, '.', ',');
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
			type: 'datetime',
			labels: {
				formatter: function() {
					return Highcharts.dateFormat('%b %y', this.value);
				}
			}
		},

		series: [
			{ name: 'Profit & Loss', color: "#757575", data: [] }
		]
	}

	//
	// Constructor
	//
	constructor(
		public activityService: ActivityService,
		public reportService: ReportService,
		public accountService: AccountService,
		private titleService: Title,
		private ledgerService: LedgerService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(this.pageTitle);

		// Load page data.
		this.refreshPage();

		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.refreshPage();
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
	// Load all the data for the page.
	//
	refreshPage() {
		// Load page data
		this.loadActivity();

		// Ledger data
		this.loadLedgerData();
	}

	//
	// Load ledger data.
	//
	loadLedgerData() {
		// Load ledger entries
		this.ledgerService.get(0, "", "", null, [], null).subscribe(res => {
			if (res.Data.length == 0) {
				this.firstRun = true;
				this.normalPage = false;
			} else {
				this.firstRun = false;
				this.normalPage = true;

				// Build charts - Need to do this after we figure out if this is first run or not.
				this.buildChart1();
				this.buildChart2();
			}
		});
	}

	//
	// Load activity
	//
	loadActivity() {
		this.activityService.getWithLimit(1, 10).subscribe(res => {
			this.activity = res;
		});
	}

	//
	// Get chart data - chart #1
	//
	buildChart1() {
		// Set date range
		let startStr = moment().subtract(4, 'months').format('YYYY-MM-01');
		let startDate = moment(startStr).toDate();
		let endDate = moment().endOf('month').toDate();

		// Ajax call to get data
		this.reportService.getPnl(startDate, endDate, "month", "asc").subscribe((res) => {
			let data1 = [];
			let data2 = [];

			for (let i = 0; i < res.length; i++) {
				data1.push({ x: res[i].Date, y: res[i].Income });
				data2.push({ x: res[i].Date, y: (res[i].Expense * -1) });
			}

			// Rebuilt the chart
			this.chartOptions1.series[0].data = data1;
			this.chartOptions1.series[1].data = data2;
			Highcharts.chart('chart1', this.chartOptions1);
		});
	}

	//
	// Get chart data - chart #2
	//
	buildChart2() {
		// Set date range
		let startStr = moment().subtract(12, 'months').format('YYYY-MM-01');
		let startDate = moment(startStr).toDate();
		let endDate = moment().endOf('month').toDate();

		// Ajax call to get data
		this.reportService.getPnl(startDate, endDate, "month", "asc").subscribe((res) => {
			let data = [];

			for (let i = 0; i < res.length; i++) {
				data.push({ x: res[i].Date, y: res[i].Profit });
			}

			// Rebuilt the chart
			this.chartOptions2.series[0].data = data;
			Highcharts.chart('chart2', this.chartOptions2);
		});
	}

	//
	// Format the message to our liking
	//
	printMessage(row: Activity) {
		// Set contact name.
		let contactName = "";
		if (row.LedgerId > 0) {
			contactName = row.Ledger.Contact.Name;

			if (contactName.length == 0) {
				contactName = `${row.Ledger.Contact.FirstName} ${row.Ledger.Contact.LastName}`;
			}
		}

		// Snapclerk message
		if (row.Action == "snapclerk") {
			if (row.SubAction == "create") {
				return row.Message.replace(row.User.FirstName, `<strong>${row.User.FirstName}</strong>`) + " To view " + `<a href="/snapclerk">click here</a>.`;
			}

			if (row.SubAction == "update") {
				return row.Message.replace(row.User.FirstName, `<strong>${row.User.FirstName}</strong>`).replace(contactName, `<a href="/ledger/${row.LedgerId}">${contactName}</a>`);
			}
		}

		// Process body
		let a = row.Message.split(" ");
		let first = a[0];
		a.shift();
		let body = a.join(" ");
		body = body.split(row.Name)[0];


		// Ledger processing
		if (row.SubAction != "delete") {
			body = `${body}<a href="/ledger/${row.LedgerId}">${row.Name}</a>`
		} else {
			body = `${body} ${row.Name}`
		}

		if (row.SubAction == "other") {
			return row.Name;
		}

		return `<strong>${first}</strong> ${body}.`;
	}
}

/* End File */
