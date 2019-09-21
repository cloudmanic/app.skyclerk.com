//
// Date: 2019-09-20
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { PipeTransform, Pipe } from '@angular/core';
import { AccountService } from '../services/account.service';

@Pipe({
	name: 'currencyFormat',
	pure: false
})

export class CurrencyFormatPipe implements PipeTransform {
	constructor(public accountService: AccountService) { }

	transform(value: number): any {
		let acct = this.accountService.getActiveAccount();
		return new Intl.NumberFormat(acct.Locale, { style: 'currency', currency: acct.Currency }).format(value);
	}
}

/* End File */
