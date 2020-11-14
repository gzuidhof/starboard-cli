/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/vfsgen"
)

func main() {
	var fs http.FileSystem = http.Dir("web/static/")

	os.MkdirAll("assets/web_static/", 0755)
	err := vfsgen.Generate(fs, vfsgen.Options{Filename: "assets/web_static/vfsdata.go", VariableName: "StaticAssets", PackageName: "web_static"})
	if err != nil {
		log.Fatalln(err)
	}

	var templateFs http.FileSystem = http.Dir("web/template/")
	os.MkdirAll("assets/web_template/", 0755)
	err = vfsgen.Generate(templateFs, vfsgen.Options{Filename: "assets/web_template/vfsdata.go", VariableName: "TemplateAssets", PackageName: "web_template"})
	if err != nil {
		log.Fatalln(err)
	}
}
