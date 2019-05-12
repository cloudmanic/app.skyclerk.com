import { Component, OnInit, Output, EventEmitter } from '@angular/core';

@Component({
	selector: 'app-contacts-assign-field',
	templateUrl: './assign-field.component.html'
})
export class AssignFieldComponent implements OnInit {
	@Output() addContactToggle = new EventEmitter<boolean>();

	showAddContact: boolean = false;

	constructor() { }

	ngOnInit() {
	}

	//
	// Toggle the add contact popover
	//
	addContactToggleClick() {
		this.addContactToggle.emit(true);
	}

}
