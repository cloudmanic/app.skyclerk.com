//
// Date: 2020-05-10
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

import * as moment from 'moment';
import { Serializable } from './serializable.model';

export class Billing implements Serializable {
	Id: number = 0;
	Subscription: string = "";
	Status: string = "en-US";
	PaymentProcessor: string = "";
	TrialExpire: Date = new Date();
	CardBrand: string = "";
	CardLast4: string = "";
	CardExpMonth: number = 0;
	CardExpYear: number = 0;
	CurrentPeriodStart: Date = new Date();
	CurrentPeriodEnd: Date = new Date();


	//
	// Json to Object.
	//
	deserialize(json: Object): this {
		this.Id = json["id"];
		this.Subscription = json["subscription"];
		this.PaymentProcessor = json["payment_processor"];
		this.Status = json["status"];
		this.TrialExpire = moment(json["trial_expire"]).toDate();
		this.CardBrand = json["card_brand"];
		this.CardLast4 = json["card_last_4"];
		this.CardExpMonth = json["card_exp_month"];
		this.CardExpYear = json["card_exp_year"];
		this.CurrentPeriodStart = moment(json["current_period_start"]).toDate();
		this.CurrentPeriodEnd = moment(json["current_period_end"]).toDate();
		return this;
	}

	//
	// Model to JS Object.
	//
	serialize(obj: Billing): Object {
		let rt = {
			id: obj.Id,
			subscription: obj.Subscription,
			status: obj.Status,
			trial_expire: obj.TrialExpire,
			card_brand: obj.CardBrand,
			card_exp_month: obj.CardExpMonth,
			card_exp_year: obj.CardExpYear,
			card_last_4: obj.CardLast4,
			current_period_start: obj.CurrentPeriodStart,
			current_period_end: obj.CurrentPeriodEnd,
			payment_processor: obj.PaymentProcessor
		}
		return rt;
	}
}

/* End File */
