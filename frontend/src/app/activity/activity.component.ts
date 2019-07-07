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
import { MeService } from '../services/me.service';

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
	constructor(public activityService: ActivityService, public meService: MeService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Load page data
		this.loadActivity();

		// Listen for account changes.
		this.meService.accountChange.takeUntil(this.destory).subscribe(() => {
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
		let a = row.Message.split(" ");
		let first = a[0];
		a.shift();
		let body = a.join(" ");
		body = body.split(row.Name)[0];

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
