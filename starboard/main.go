/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

//go:generate go run scripts/download_runtime/main.go starboard-notebook 0.10.1 web/static/vendor/
//go:generate go run scripts/download_runtime/main.go starboard-wrap 0.2.5 web/static/vendor/

import (
	"github.com/gzuidhof/starboard-cli/starboard/cmd"
)

func main() {
	cmd.Execute()
}
