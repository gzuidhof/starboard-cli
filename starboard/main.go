/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

//go:generate go run web/assets_generate.go

import (
	"github.com/gzuidhof/starboard-cli/starboard/cmd"
)

func main() {
	cmd.Execute()
}
