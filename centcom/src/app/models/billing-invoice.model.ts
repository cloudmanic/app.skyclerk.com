//
// Date: 2020-05-10
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Serializable } from './serializable.model';

export class BillingInvoice implements Serializable {
	Date: Date = new Date();
	Amount: number = 0.00
	Transaction: string = "";
	PaymentMethod: string = "";
	InvoiceURL: string = "";

	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Date = moment(json["date"]).toDate();
		this.Amount = json["amount"];
		this.Transaction = json["transaction"];
		this.PaymentMethod = json["payment_method"];
		this.InvoiceURL = json["invoice_url"];
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: BillingInvoice): Object {
		let rt = {
			date: obj.Date,
			amount: obj.Amount,
			payment_method: obj.PaymentMethod,
			transaction: obj.Transaction,
			invoice_url: obj.InvoiceURL,
		}
		return rt;
	}
}

/* End File */
