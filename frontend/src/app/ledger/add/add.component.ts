import { Component, OnInit, Input } from '@angular/core';

@Component({
	selector: 'app-ledger-add',
	templateUrl: './add.component.html'
})
export class AddComponent implements OnInit {
	@Input() type: string = "income";

	showAddContact: boolean = false;

	constructor() { }

	ngOnInit() {
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
