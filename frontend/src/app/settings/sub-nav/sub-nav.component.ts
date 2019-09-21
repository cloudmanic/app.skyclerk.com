import { Component, OnInit, Input } from '@angular/core';
import { Account } from 'src/app/models/account.model';
import { MeService } from 'src/app/services/me.service';
import { AccountService } from 'src/app/services/account.service';
import { Me } from 'src/app/models/me.model';

@Component({
	selector: 'app-settings-sub-nav',
	templateUrl: './sub-nav.component.html'
})
export class SubNavComponent implements OnInit {
	@Input() current: string = "";

	me: Me = new Me();
	account: Account = new Account();
	isOwner: boolean = false;

	//
	// Construct.
	//
	constructor(public meService: MeService, public accountService: AccountService) { }

	//
	// ngOnInit
	//
	ngOnInit() {
		// Get Me.
		this.meService.get().subscribe(res => {
			this.me = res;
			this.setOwner();
		});

		// Get the active account.
		this.accountService.getAccount().subscribe(res => {
			this.account = res;
			this.setOwner();
		});
	}

	//
	// See if we are an owner
	//
	setOwner() {
		if (this.me.Id == this.account.OwnerId) {
			this.isOwner = true;
		} else {
			this.isOwner = false;
		}
	}
}

/* End File */
