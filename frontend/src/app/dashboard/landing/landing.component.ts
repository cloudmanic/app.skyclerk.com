//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { ActivityService } from 'src/app/services/activity.service';
import { Activity } from 'src/app/models/activity.model';

@Component({
	selector: 'app-landing',
	templateUrl: './landing.component.html'
})

export class LandingComponent implements OnInit {
	activity: Activity[] = [];
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
		this.activityService.get(1, 25).subscribe(res => {
			this.activity = res;
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

		if (row.SubAction != "delete") {
			body = `${body}<a href="/ledger/${row.LedgerId}">${row.Name}</a>`
		} else {
			body = `${body} ${row.Name}`
		}

		return `<strong>${first}</strong> ${body}.`;
	}
}

/* End File */
