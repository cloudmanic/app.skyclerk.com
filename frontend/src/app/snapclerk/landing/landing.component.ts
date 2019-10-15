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
import { File as FileModel } from 'src/app/models/file.model';
import { Title } from '@angular/platform-browser';
import { environment } from 'src/environments/environment';
import { AccountService } from 'src/app/services/account.service';

const pageTitle: string = environment.title_prefix + "Snap!Clerk";

@Component({
	selector: 'snapclerk-app-landing',
	templateUrl: './landing.component.html'
})

export class LandingComponent implements OnInit {
	firstRun: boolean = false;
	showSnapClerk: boolean = false;
	snapclerks: SnapClerkResponse = new SnapClerkResponse(false, 0, 50, 0, []);
	page: number = 1;
	usage: number = 0;
	destory: Subject<boolean> = new Subject<boolean>();

	//
	// Constructor
	//
	constructor(public snapClerkService: SnapClerkService, public meService: MeService, public titleService: Title, public accountService: AccountService) { }

	//
	// NgOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);

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
	// refreshPage - load the data for the page.
	//
	refreshPage() {
		// Get the list of snapclerks
		this.snapClerkService.get(this.page, 25, "SnapClerkId", "DESC").subscribe(res => {
			// Logic to see if this is the first time we are calling this.
			if ((!this.showSnapClerk) && (!this.firstRun)) {
				if (res.Data.length == 0) {
					this.firstRun = true;
				} else {
					this.showSnapClerk = true;
				}
			}

			this.snapclerks = res;
		});

		// Get the snapclerk usage
		this.snapClerkService.getUsage().subscribe(res => {
			this.usage = res;
		});
	}

	//
	// Close the first run.
	//
	closeFirstRun() {
		this.firstRun = false;
		this.showSnapClerk = true;
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

	//
	// We call this when a file has been uploaded
	//
	onAddFile(f: FileModel) {
		// Send file to server.
		this.snapClerkService.create(f.Id).subscribe(_res => {
			this.refreshPage();
		});
	}
}

/* End File */
