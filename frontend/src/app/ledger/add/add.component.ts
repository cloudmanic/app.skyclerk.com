//
// Date: 2019-05-05
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { Component, OnInit, Input } from '@angular/core';
import { Contact } from 'src/app/models/contact.model';
import { Ledger } from 'src/app/models/ledger.model';
import { Category } from 'src/app/models/category.model';
import { Label } from 'src/app/models/label.model';
import { File as FileModel } from 'src/app/models/file.model';
import { Router } from '@angular/router';

@Component({
	selector: 'app-ledger-add',
	templateUrl: './add.component.html'
})
export class AddComponent implements OnInit {
	@Input() type: string = "income";

	ledger: Ledger = new Ledger();
	showAddLabel: boolean = false;
	showAddContact: boolean = false;
	showAddCategory: boolean = false;

	//
	// Constructor
	//
	constructor(public router: Router) { }

	//
	// ngOnInit
	//
	ngOnInit() { }

	//
	// Save the new ledger entry.
	//
	save() {
		console.log(this.ledger);
		this.router.navigate(['/ledger']);
	}

	//
	// We call this on assigning a contact.
	//
	assignContact(contact: Contact) {
		this.showAddContact = false;
		this.ledger.Contact = contact;
	}

	//
	// We call this whenever someone checks or unchecks a label
	//
	onLabelsChange(lbs: Label[]) {
		this.ledger.Labels = lbs;
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
	// We call this after assigning a label.
	//
	assignLabel(lb: Label) {
		this.showAddLabel = false;
		this.ledger.Labels.push(lb);

		// Hack to make sure things update.
		this.onLabelsChange(this.ledger.Labels.slice());
	}

	//
	// We call this after assigning a category.
	//
	assignCategory(category: Category) {
		this.showAddCategory = false;
		this.ledger.Category = category;
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
	// Add contact click
	//
	addContactToggle() {
		if (this.showAddContact) {
			this.showAddContact = false;
		} else {
			this.showAddContact = true;
		}
	}

}

/* End File */
