/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"net/http"
	"strings"

	"github.com/gzuidhof/starboard-cli/starboard/assets/web_static"
	"github.com/gzuidhof/starboard-cli/starboard/assets/web_templates"
	"github.com/spf13/viper"
)

const dev = true

type serveFS struct {
	static    http.FileSystem
	templates http.FileSystem
}

func getFileSystems() serveFS {
	var static http.FileSystem
	var templates http.FileSystem

	if viper.GetString("static_folder") != "" {
		static = http.Dir(viper.GetString("static_folder"))
	} else {
		static = web_static.StaticAssets
	}

	if viper.GetString("templates_folder") != "" {
		templates = http.Dir(viper.GetString("templates_folder"))
	} else {
		templates = web_templates.TemplateAssets
	}

	return serveFS{
		static:    static,
		templates: templates,
	}
}

func isProbablyNotebookFilename(name string) bool {
	return strings.HasSuffix(name, ".nb") || strings.HasSuffix(name, ".sb") || strings.HasSuffix(name, ".sbnb")
}
