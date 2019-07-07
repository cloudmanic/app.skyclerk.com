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
import { Contact } from 'src/app/models/contact.model';
import { Label } from 'src/app/models/label.model';
import { Category } from 'src/app/models/category.model';
import { File as FileModel } from 'src/app/models/file.model';
import { MeService } from 'src/app/services/me.service';
import { Subject } from 'rxjs';

@Component({
	selector: 'app-ledger-edit',
	templateUrl: './edit.component.html'
})
export class EditComponent implements OnInit {
	errors: any = [];
	type: string = "income";
	ledger: Ledger = new Ledger();
	destory: Subject<boolean> = new Subject<boolean>();

	showAddLabel: boolean = false;
	showAddContact: boolean = false;
	showAddCategory: boolean = false;

	//
	// Constructor
	//
	constructor(public ledgerService: LedgerService, public route: ActivatedRoute, public router: Router, public meService: MeService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Is this an edit action?
		let ledgerId = this.route.snapshot.params['id'];

		// Get the ledger based on the id we passed in.
		this.ledgerService.getById(ledgerId).subscribe(res => {
			this.ledger = res;

			// Set type.
			if (this.ledger.Amount < 0) {
				this.type = "expense";
			}
		});

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
	// Update the ledger entry.
	//
	save() {
		// Clear errors
		this.errors = [];

		// Is this an expense
		if (this.type == "expense") {
			this.ledger.Amount = (this.ledger.Amount * -1);
		}

		// Send this ledger to BE.
		this.ledgerService.update(this.ledger).subscribe(
			// Sucesss
			() => {
				this.router.navigate([`/ledger/${this.ledger.Id}`]);
			},

			// Error
			err => {
				// Reset negative value
				if (this.type == "expense") {
					this.ledger.Amount = (this.ledger.Amount * -1);
				}

				// Show errors
				this.errors = err.error.errors;
			}
		);
	}

	//
	// We call this after assigning a category.
	//
	assignCategory(category: Category) {
		this.showAddCategory = false;
		this.ledger.Category = category;
	}

	//
	// We call this after assigning a label.
	//
	assignLabel(lb: Label) {
		this.showAddLabel = false;
		this.ledger.Labels.push(lb);

		// Hack to make sure things update.
		this.onLabelsChange(this.ledger.Labels.slice());
	}

	//
	// We call this whenever someone checks or unchecks a label
	//
	onLabelsChange(lbs: Label[]) {
		this.ledger.Labels = lbs;
	}

	//
	// We call this on assigning a contact.
	//
	assignContact(contact: Contact) {
		this.showAddContact = false;
		this.ledger.Contact = contact;
	}

	//
	// Add contact click
	//
	addContactToggle() {
		if (this.showAddContact) {
			this.showAddContact = false;
		} else {
			this.showAddContact = true;
		}
	}

	//
	// Add label click
	//
	addLabelToggle() {
		if (this.showAddLabel) {
			this.showAddLabel = false;
		} else {
			this.showAddLabel = true;
		}
	}

	//
	// Add category click
	//
	addCategoryToggle() {
		if (this.showAddCategory) {
			this.showAddCategory = false;
		} else {
			this.showAddCategory = true;
		}
	}

	//
	// We call this when a file is added to the ledger.
	//
	onAddFile(f: FileModel) {
		this.ledger.Files.push(f);
	}

	//
	// Call this when we delete a file.
	//
	onDeleteFile(f: FileModel) {
		for (let i = 0; i < this.ledger.Files.length; i++) {
			if (f.Id == this.ledger.Files[i].Id) {
				this.ledger.Files.splice(i, 1);
			}
		}
	}

	//
	// Check to see if we have errors for a field.
	//
	hasError(field: string) {
		if (this.errors[field]) {
			return this.errors[field];
		}

		return "";
	}
}

/* End File */
