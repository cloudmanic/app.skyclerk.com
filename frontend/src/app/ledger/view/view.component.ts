//
// Date: 2019-05-20
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import 'rxjs/add/operator/takeUntil';
import { Component, OnInit } from '@angular/core';
import { LedgerService } from 'src/app/services/ledger.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Ledger } from 'src/app/models/ledger.model';
import { Activity } from 'src/app/models/activity.model';
import { ActivityService } from 'src/app/services/activity.service';
import { MeService } from 'src/app/services/me.service';
import { Subject } from 'rxjs';
import { Title } from '@angular/platform-browser';
import { environment } from 'src/environments/environment';

const pageTitle: string = environment.title_prefix + "Ledger View";

@Component({
	selector: 'app-ledger-view',
	templateUrl: './view.component.html'
})

export class ViewComponent implements OnInit {
	activity: Activity[] = [];
	ledger: Ledger = new Ledger();
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public ledgerService: LedgerService, public route: ActivatedRoute, public router: Router, public activityService: ActivityService, public meService: MeService, private titleService: Title) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

		// Is this an edit action?
		let ledgerId = this.route.snapshot.params['id'];

		// Get the ledger based on the id we passed in.
		this.loadLedgerEntry(ledgerId);

		// Load activity for this ledger entry.
		this.loadActivity(ledgerId);

		// Listen for account changes.
		this.meService.accountChange.takeUntil(this.destory).subscribe(() => {
			this.router.navigate([`/`]);
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
	// Load Ledger entry.
	//
	loadLedgerEntry(ledgerId: number) {
		this.ledgerService.getById(ledgerId).subscribe(res => {
			this.ledger = res;
		});
	}

	//
	// Load activity
	//
	loadActivity(ledgerId: number) {
		this.activityService.getByLedgerId(ledgerId).subscribe(res => {
			this.activity = res;

			console.log(this.activity);
		});
	}

	//
	// Delete ledger
	//
	deleteLedger() {
		let c = confirm("Are you sure you want to delete this ledger entry?");

		if (!c) {
			return;
		}

		// Send delete request.
		this.ledgerService.delete(this.ledger).subscribe(() => {
			this.router.navigate(['/ledger']);
		});
	}
}

/* End File */
