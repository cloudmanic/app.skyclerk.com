//
// Date: 2019-07-7
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { SnapClerkService, SnapClerkResponse } from 'src/app/services/snapckerk.service';
import { Subject } from 'rxjs';
import { MeService } from 'src/app/services/me.service';

@Component({
	selector: 'snapclerk-app-landing',
	templateUrl: './landing.component.html'
})

export class LandingComponent implements OnInit {
	snapclerks: SnapClerkResponse = new SnapClerkResponse(false, 0, 50, 0, []);
	page: number = 1;
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public snapClerkService: SnapClerkService, public meService: MeService) { }

	//
	// NgOnInit
	//
	ngOnInit() {
		// Load page data.
		this.refreshPage();

		// Listen for account changes.
		this.meService.accountChange.takeUntil(this.destory).subscribe(() => {
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
	// refreshPage - load the data for the page.
	//
	refreshPage() {
		// Get the list of snapclerks
		this.snapClerkService.get(this.page, "SnapClerkId", "DESC").subscribe(res => {
			this.snapclerks = res;
			console.log(this.snapclerks);
		});
	}

	//
	// Return the page list for ledger
	//
	getPageRange() {
		let pages = [];

		if (this.snapclerks.Data.length == 0) {
			return [1];
		}

		let pageCount = Math.ceil(this.snapclerks.NoLimitCount / this.snapclerks.Limit);

		for (let i = 1; i <= pageCount; i++) {
			pages.push(i);
		}

		return pages;
	}

	//
	// Paging select change
	//
	doPageSelectChange() {
		this.refreshPage();
	}

	//
	// Paging next click
	//
	doNextClick() {
		this.page++;
		this.refreshPage();
	}

	//
	// Paging prev click
	//
	doPrevClick() {
		this.page--;
		this.refreshPage();
	}
}

/* End File */
