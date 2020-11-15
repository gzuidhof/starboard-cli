/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"strings"
)

type BreadCrumb struct {
	Name string
	/**
	* Path without leading "/"
	 */
	Path     string
	IsFolder bool
}

func makeBreadCrumbs(path string, lastIsFolder bool) []BreadCrumb {
	path = strings.Trim(path, "/")
	if path == "" {
		return []BreadCrumb{}
	}
	parts := strings.Split(path, "/")

	bc := make([]BreadCrumb, len(parts))

	for i := 0; i < len(parts); i++ {
		bc[i] = BreadCrumb{
			Name:     parts[i],
			IsFolder: i < len(parts)-1 || lastIsFolder,
			Path:     strings.Join(parts[:i+1], "/"),
		}
	}

	return bc
}
