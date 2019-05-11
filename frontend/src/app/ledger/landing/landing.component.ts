//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit } from '@angular/core';
import { LedgerService, LedgerResponse, LedgerSummaryResponse, LedgerPnlSummary } from 'src/app/services/ledger.service';
import { MeService } from 'src/app/services/me.service';
import { Category } from 'src/app/models/category.model';
import { Label } from 'src/app/models/label.model';

@Component({
	selector: 'app-landing',
	templateUrl: './landing.component.html'
})

export class LandingComponent implements OnInit {
	page: number = 1;
	plSummary: LedgerPnlSummary = new LedgerPnlSummary(0, 0, 0);
	pageRangeSelect: number = 1;
	ledgerSummary: LedgerSummaryResponse = new LedgerSummaryResponse([], [], []);
	ledgers: LedgerResponse = new LedgerResponse(false, 0, 25, 0, []);

	// Active filters
	type: string = "";
	search: string = "";
	activeYear: number = null;
	activeLabels: Label[] = [];
	activeCategory: Category = null;


	//
	// Construct
	//
	constructor(public ledgerService: LedgerService, public meService: MeService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Load starting data
		this.refreshLedger();

		// Listen for account changes.
		this.meService.accountChange.subscribe(() => {
			this.refreshLedger();
		});
	}

	//
	// Refresh ledger.
	//
	refreshLedger() {
		this.loadLedgerData();
		this.getLedgerSummary();
		this.getLedgerPnlSummary();
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
		this.refreshLedger();
	}

	//
	// Load ledger data.
	//
	loadLedgerData() {
		// Load ledger entries
		this.ledgerService.get(this.page, this.type, this.search, this.activeCategory, this.activeLabels, this.activeYear).subscribe(res => {
			this.ledgers = res;
			this.pageRangeSelect = this.page;
		});
	}

	//
	// getLedgerPnlSummary data.
	//
	getLedgerPnlSummary() {
		this.ledgerService.getLedgerPnlSummary(this.type, this.search, this.activeCategory, this.activeLabels, this.activeYear).subscribe(res => {
			this.plSummary = res;
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
		this.refreshLedger();
	}

	//
	// Paging next click
	//
	doNextClick() {
		this.page++;
		this.refreshLedger();
	}

	//
	// Paging prev click
	//
	doPrevClick() {
		this.page--;
		this.refreshLedger();
	}

	// ------------- Filter Functions --------------- //

	//
	// Do we show filters
	//
	doWeShowFilters() {
		if (this.activeYear) {
			return true;
		}

		if (this.activeCategory) {
			return true;
		}

		if (this.activeLabels.length > 0) {
			return true;
		}

		return false;
	}

	//
	// Add Label fitter
	//
	addLabelFilter(lb: Label) {
		this.page = 1;
		this.activeLabels.push(lb);
		this.refreshLedger();
		this.getLedgerPnlSummary();
	}

	//
	// Remove Label fitter
	//
	removeLabelFilter(lb: Label) {
		this.page = 1;
		this.activeLabels = this.activeLabels.filter(obj => obj !== lb);
		this.refreshLedger();
		this.getLedgerPnlSummary();
	}

	//
	// Set year filter
	//
	setYearFilter(year: number) {
		this.page = 1;
		this.activeYear = year;
		this.refreshLedger();
		this.getLedgerPnlSummary();
	}

	//
	// Set category filter
	//
	setCategoryFilter(cat: Category) {
		this.page = 1;
		this.activeCategory = cat;
		this.refreshLedger();
		this.getLedgerPnlSummary();
	}

	//
	// Clear all sidebar filters
	//
	clearAllSidebarFilters() {
		this.page = 1;
		this.activeYear = null;
		this.activeLabels = [];
		this.activeCategory = null;
		this.refreshLedger();
		this.getLedgerPnlSummary();
	}

	//
	// Do search.
	//
	doSearch() {
		this.page = 1;
		this.refreshLedger();
		this.getLedgerPnlSummary();
	}

	//
	// Change the type we are filtring by
	//
	doTypeClick(type: string) {
		this.type = type;
		this.page = 1;
		this.refreshLedger();
		this.getLedgerPnlSummary();
	}

	//
	// Return the count of a year
	//
	filterGetYearCount(year: number) {
		// Match up the label
		for (let i = 0; i < this.ledgerSummary.Years.length; i++) {
			if (this.ledgerSummary.Years[i].Year == year) {
				return this.ledgerSummary.Years[i].Count;
			}
		}
		// Should never get here.
		return 0;
	}

	//
	// Return the count of a label
	//
	filterGetLabelCount(lb: Label) {
		// Match up the label
		for (let i = 0; i < this.ledgerSummary.Labels.length; i++) {
			if (this.ledgerSummary.Labels[i].Id == lb.Id) {
				return this.ledgerSummary.Labels[i].Count;
			}
		}
		// Should never get here.
		return 0;
	}

	//
	// Filter get category count
	//
	filterGetCategoryCount(cat: Category) {
		// Match up the category
		for (let i = 0; i < this.ledgerSummary.Categories.length; i++) {
			if (this.ledgerSummary.Categories[i].Id == cat.Id) {
				return this.ledgerSummary.Categories[i].Count;
			}
		}
		// Should never get here.
		return 0;
	}
}

/* End File */
