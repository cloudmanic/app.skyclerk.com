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
	snapclerks: SnapClerkResponse;
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public snapClerkService: SnapClerkService, public meService: MeService) { }

	//
	// NgOnInit
	//
	ngOnInit() {
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
		this.snapClerkService.get(1, "desc", "SnapClerkId").subscribe(res => {
			console.log(res);
			this.snapclerks = res;
		});
	}
}

/* End File */
