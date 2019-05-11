import { Component, OnInit, Input } from '@angular/core';

@Component({
	selector: 'app-contacts-add-edit-pop-out',
	templateUrl: './add-edit-pop-out.component.html'
})
export class AddEditPopOutComponent implements OnInit {
	@Input() show: boolean = false;

	constructor() { }

	ngOnInit() {
	}

	//
	// Cancel popup.
	//
	cancelPopUp() {
		this.show = false;
	}
}
