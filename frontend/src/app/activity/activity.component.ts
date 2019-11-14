//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Component, OnInit } from '@angular/core';
import { Activity } from '../models/activity.model';
import { ActivityService, ActivityResponse } from '../services/activity.service';
import { Subject } from 'rxjs';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';
import { AccountService } from '../services/account.service';

const pageTitle: string = environment.title_prefix + "Activity";

@Component({
	selector: 'app-activity',
	templateUrl: './activity.component.html'
})

export class ActivityComponent implements OnInit {
	page: number = 1;
	activity: any;
	activityResponse: ActivityResponse = new ActivityResponse(false, 0, 25, 0, []);
	activityKeys: String[];
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public activityService: ActivityService, public accountService: AccountService, public titleService: Title) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Load page data
		this.loadActivity();

		// Listen for account changes.
		this.accountService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.loadActivity();
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
	// Load activity
	//
	loadActivity() {
		this.activityService.get(this.page, 25, "id", "DESC").subscribe(res => {
			this.activityResponse = res;

			// Build Grouping
			this.activity = {}
			this.activityKeys = [];

			for (let i = 0; i < res.Data.length; i++) {
				let ix = moment(res.Data[i].CreatedAt).format("YYYY-MM-DD");

				if (typeof this.activity[ix] == "undefined") {
					this.activity[ix] = [];
					this.activityKeys.push(ix);
				}

				this.activity[ix].push(res.Data[i]);
			}
		});
	}

	//
	// Return the page list for ledger
	//
	getPageRange() {
		let pages = [];

		if (this.activityResponse.Data.length == 0) {
			return [1];
		}

		let pageCount = Math.ceil(this.activityResponse.NoLimitCount / this.activityResponse.Limit);

		for (let i = 1; i <= pageCount; i++) {
			pages.push(i);
		}

		return pages;
	}

	//
	// Paging select change
	//
	doPageSelectChange() {
		this.loadActivity();
	}

	//
	// Paging next click
	//
	doNextClick() {
		this.page++;
		this.loadActivity();
	}

	//
	// Paging prev click
	//
	doPrevClick() {
		this.page--;
		this.loadActivity();
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
