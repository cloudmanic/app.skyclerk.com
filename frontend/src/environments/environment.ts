// This file can be replaced during build by using the `fileReplacements` array.
// `ng build --prod` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

export const environment = {
	production: false,
	version: "1.0.0",
	client_id: "XL8TeRGBdsUvvM3",
	app_server: "http://127.0.0.1:9090",
	mixpanel_key: "",
	stripe_pub_key: "pk_test_AZaPyHOuCrE4AMf9ZzQts8IH00RtU5m2Cd",
	//mixpanel_key: "5974bde1c90a5606add2695a17b2db10",
	title_prefix: "Skyclerk | "
};

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
import 'zone.js/dist/zone-error';  // Included with Angular CLI.
