//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { LedgerService, LedgerResponse, LedgerSummaryResponse } from 'src/app/services/ledger.service';
import { MeService } from 'src/app/services/me.service';

@Component({
	selector: 'app-landing',
	templateUrl: './landing.component.html'
})

export class LandingComponent implements OnInit {
	page: number = 1;
	type: string = "";
	search: string = "";
	pageRangeSelect: number = 1;
	ledgerSummary: LedgerSummaryResponse = new LedgerSummaryResponse([], [], []);
	ledgers: LedgerResponse = new LedgerResponse(false, 0, 50, 0, []);

	//
	// Construct
	//
	constructor(public ledgerService: LedgerService, public meService: MeService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Load starting data
		this.loadPageData();

		// Listen for account changes.
		this.meService.accountChange.subscribe(() => {
			this.loadPageData();
		});
	}

	//
	// Load page data
	//
	loadPageData() {
		// Reset filters, and such.
		this.page = 1;
		this.type = "";
		this.search = "";
		this.pageRangeSelect = 1;

		// Load ledger data
		this.loadLedgerData();
	}

	//
	// Load ledger data.
	//
	loadLedgerData() {
		// Load ledger summary
		this.getLedgerSummary();

		// Load ledger entries
		this.ledgerService.get(this.page, this.type, this.search).subscribe(res => {
			this.ledgers = res;
			this.pageRangeSelect = this.page;
		});
	}

	//
	// Get ledger summary
	//
	getLedgerSummary() {
		this.ledgerService.getLedgerSummary(this.type).subscribe(res => {
			this.ledgerSummary = res;
		});
	}

	//
	// Do search.
	//
	doSearch() {
		this.page = 1;
		this.loadLedgerData();
	}

	//
	// Change the type we are filtring by
	//
	doTypeClick(type: string) {
		this.type = type;
		this.page = 1;
		this.loadLedgerData()
	}

	//
	// Return the page list for ledger
	//
	getPageRange() {
		let pages = [];

		if (this.ledgers.Data.length == 0) {
			return [1];
		}

		let pageCount = Math.ceil(this.ledgers.NoLimitCount / this.ledgers.Limit);

		for (let i = 1; i <= pageCount; i++) {
			pages.push(i);
		}

		return pages;
	}

	//
	// Paging select change
	//
	doPageSelectChange() {
		this.page = this.pageRangeSelect;
		this.loadLedgerData();
	}

	//
	// Paging next click
	//
	doNextClick() {
		this.page++;
		this.loadLedgerData();
	}

	//
	// Paging prev click
	//
	doPrevClick() {
		this.page--;
		this.loadLedgerData();
	}
}

/* End File */
