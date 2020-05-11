//
// Date: 2020-05-10
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AccountService } from 'src/app/services/account.service';
import { Subject } from 'rxjs';

@Component({
	selector: 'app-plans',
	templateUrl: './plans.component.html'
})
export class PlansComponent implements OnInit {
	back: string = "";
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor.
	//
	constructor(public route: ActivatedRoute, public accountService: AccountService, public router: Router) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set plan on load.
		this.back = this.route.snapshot.queryParamMap.get("back");

		// When plan changes
		this.route.queryParamMap.takeUntil(this.destory).subscribe(queryParams => {
			this.back = queryParams.get("back");
		});
	}

	//
	// OnDestroy
	//
	ngOnDestroy() {
		this.destory.next();
		this.destory.complete();
	}
}

/* End File */
