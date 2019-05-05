import { Component, OnInit, Input } from '@angular/core';

@Component({
	selector: 'app-settings-sub-nav',
	templateUrl: './sub-nav.component.html'
})
export class SubNavComponent implements OnInit {
	@Input() current: string = "";

	constructor() { }

	ngOnInit() {
	}

}
