/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var fs http.FileSystem = http.Dir("web/static/")
	err := vfsgen.Generate(fs, vfsgen.Options{Filename: "assets/web_static_vfsdata.go", VariableName: "WebStaticAssets", PackageName: "assets"})
	if err != nil {
		log.Fatalln(err)
	}
}
