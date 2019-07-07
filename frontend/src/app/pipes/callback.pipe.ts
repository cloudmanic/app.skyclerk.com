//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

import { PipeTransform, Pipe } from '@angular/core';

@Pipe({
	name: 'callback',
	pure: false
})

export class CallbackPipe implements PipeTransform {
	transform(items: any[], callback: (item: any) => boolean): any {
		if (!items || !callback) {
			return items;
		}
		return items.filter(item => callback(item));
	}
}

/* End File */
