import { Component, OnInit, Input } from '@angular/core';

@Component({
	selector: 'app-ledger-add-edit',
	templateUrl: './add-edit.component.html'
})
export class AddEditComponent implements OnInit {
	@Input() type: string = "income";

	showAddContact: boolean = false;


	constructor() { }

	ngOnInit() {
	}

	//
	// Add contact click
	//
	addContactClick() {
		if (this.showAddContact) {
			this.showAddContact = false;
		} else {
			this.showAddContact = true;
		}
	}

}

/* End File */
