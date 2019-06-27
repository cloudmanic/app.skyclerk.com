//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { Activity } from '../models/activity.model';
import { ActivityService } from '../services/activity.service';

@Component({
	selector: 'app-activity',
	templateUrl: './activity.component.html'
})

export class ActivityComponent implements OnInit {
	activity: any;
	activityKeys: Number[];

	//
	// Constructor
	//
	constructor(public activityService: ActivityService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Load page data
		this.loadActivity();
	}

	//
	// Load activity
	//
	loadActivity() {
		this.activityService.getByGroup(1, 25).subscribe(res => {
			this.activity = res;

			let t = [];
			for (let row in this.activity) {
				t.push(row);
			}

			this.activityKeys = t.slice().reverse();
		});
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
		body = `${body}<a href="/ledger/${row.LedgerId}">${row.Name}</a>`
		return `<strong>${first}</strong> ${body}.`;
	}
}

/* End File */
