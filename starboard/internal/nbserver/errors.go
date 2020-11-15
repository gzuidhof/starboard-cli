/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"net/http"
	"os"
)

// Adapted from Go http sourcecode
func toHTTPError(err error) (msg string, httpStatus int) {

	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}

	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}

	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
