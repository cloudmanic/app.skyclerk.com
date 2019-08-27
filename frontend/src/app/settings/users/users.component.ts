import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';
import { Title } from '@angular/platform-browser';

const pageTitle: string = environment.title_prefix + "Settings Users";

@Component({
	selector: 'app-users',
	templateUrl: './users.component.html'
})

export class UsersComponent implements OnInit {
	//
	// Construct.
	//
	constructor(public titleService: Title) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Set page title.
		this.titleService.setTitle(pageTitle);
	}
}

/* End File */
