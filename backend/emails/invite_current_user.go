//
// Date: 2019-07-02
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package emails

import "fmt"

//
// GetInviteCurrentUserText will set text
//
func GetInviteCurrentUserText(name string, accountName string, url string) string {
	return fmt.Sprintf(`
		%s added you to the %s Skyclerk account.
		Click Here: %s
	`, name, accountName, url)
}

//
// GetInviteCurrentUserHtml will set html
//
func GetInviteCurrentUserHTML(name string, accountName string, url string) string {
	return fmt.Sprintf(`
		<p><b>%s</b> added you to the <b>%s</b> Skyclerk account.</p>
		<p><a href="%s">Click here to access your account.</a></p>
	`, name, accountName, url)

}

/* End File */
