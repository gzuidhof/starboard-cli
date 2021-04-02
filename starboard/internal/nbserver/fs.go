/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"strings"

	"github.com/gzuidhof/starboard-cli/starboard/web/static"
	"github.com/gzuidhof/starboard-cli/starboard/web/templates"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const dev = true

type serveFS struct {
	static    afero.Fs
	templates afero.Fs
}

func getFileSystems() serveFS {
	var staticFS afero.Fs
	var templatesFS afero.Fs

	if viper.GetString("static_folder") != "" {
		staticFS = afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("static_folder")))
	} else {
		staticFS = afero.FromIOFS{static.FS}
	}

	if viper.GetString("templates_folder") != "" {
		templatesFS = afero.NewReadOnlyFs(afero.NewBasePathFs(afero.NewOsFs(), viper.GetString("templates_folder")))
	} else {
		templatesFS = afero.FromIOFS{templates.FS}
	}

	return serveFS{
		static:    staticFS,
		templates: templatesFS,
	}
}

func isProbablyNotebookFilename(name string) bool {
	return strings.HasSuffix(name, ".nb") || strings.HasSuffix(name, ".sb") || strings.HasSuffix(name, ".sbnb")
}
