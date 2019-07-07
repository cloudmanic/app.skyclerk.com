//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Component, OnInit } from '@angular/core';
import { Activity } from '../models/activity.model';
import { ActivityService } from '../services/activity.service';
import { Subject } from 'rxjs';
import { MeService } from '../services/me.service';

@Component({
	selector: 'app-activity',
	templateUrl: './activity.component.html'
})

export class ActivityComponent implements OnInit {
	activity: any;
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
		this.activityService.get(1, 25).subscribe(res => {
			// Build Grouping
			this.activity = {}
			this.activityKeys = [];

			for (let i = 0; i < res.length; i++) {
				let ix = moment(res[i].CreatedAt).format("YYYY-MM-DD");

				if (typeof this.activity[ix] == "undefined") {
					this.activity[ix] = [];
					this.activityKeys.push(ix);
				}

				this.activity[ix].push(res[i]);
			}
		});
	}

	//
	// Format the message to our liking
	//
	printMessage(row: Activity) {

		console.log(row);

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
