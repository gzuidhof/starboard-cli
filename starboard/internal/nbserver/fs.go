/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"net/http"
	"strings"

	"github.com/gzuidhof/starboard-cli/starboard/web/static"
	"github.com/gzuidhof/starboard-cli/starboard/web/templates"
	"github.com/spf13/viper"
)

const dev = true

type serveFS struct {
	static    http.FileSystem
	templates http.FileSystem
}

func getFileSystems() serveFS {
	var staticFS http.FileSystem
	var templatesFS http.FileSystem

	if viper.GetString("static_folder") != "" {
		staticFS = http.Dir(viper.GetString("static_folder"))
	} else {
		staticFS = http.FS(static.FS)
	}

	if viper.GetString("templates_folder") != "" {
		templatesFS = http.Dir(viper.GetString("templates_folder"))
	} else {
		templatesFS = http.FS(templates.FS)
	}

	return serveFS{
		static:    staticFS,
		templates: templatesFS,
	}
}

func isProbablyNotebookFilename(name string) bool {
	return strings.HasSuffix(name, ".nb") || strings.HasSuffix(name, ".sb") || strings.HasSuffix(name, ".sbnb")
}
