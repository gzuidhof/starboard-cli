/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gzuidhof/starboard-cli/starboard/internal/fs/stripprefix"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var indexTemplate *template.Template
var browseTemplate *template.Template
var editorTemplate *template.Template

func loadTemplates(fs *afero.HttpFs) {

	t, err := vfstemplate.ParseGlob(fs, nil, "*.tmpl")
	if err != nil {
		log.Fatalf("Failed to parse templates in vfs: %v", err)
	}
	indexTemplate = t.Lookup("index.html.tmpl")
	browseTemplate = t.Lookup("browse.html.tmpl")
	editorTemplate = t.Lookup("editor.html.tmpl")
}

func CreateServer(serveFolderAbs string, serveFS serveFS, portPrimary string, portSecondary string) {
	app := fiber.New(fiber.Config{CaseSensitive: true, DisableStartupMessage: true})

	app.Use(recover.New())
	app.Use(logger.New())

	app.Use("/static/*", filesystem.New(filesystem.Config{
		Root: afero.NewHttpFs(stripprefix.New("/static/", serveFS.static)),
	}))

	writeFileSystem := afero.NewBasePathFs(afero.NewOsFs(), serveFolderAbs).(*afero.BasePathFs)
	app.Get(defaultBrowseEndpoint+"*", adaptor.HTTPHandler(&browseHandler{
		root: http.Dir(serveFolderAbs),
	}))

	app.All(defaultNotebookEndpoint+"*", adaptor.HTTPHandler(&notebookHandler{
		root:        http.Dir(serveFolderAbs),
		iframeHost:  "http://localhost:" + portSecondary,
		writeFS:     writeFileSystem,
		serveFolder: serveFolderAbs,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		if isProbablyNotebookFilename(serveFolderAbs) {
			c.Redirect("/nb/")
		} else {
			c.Redirect("/browse/")
		}

		return nil
	})

	done := make(chan bool)
	go func() {
		log.Fatal(app.Listen(":" + portPrimary))
	}()
	go func() {
		log.Fatal(app.Listen(":" + portSecondary))
	}()
	log.Printf("Listening on :%v (and :%s for sandboxing)\nhttp://localhost:%v", portPrimary, portSecondary, portPrimary)

	<-done
}

func Start(servePath string) {
	port := viper.GetString("port")
	portSecondary := viper.GetString("port_secondary")
	serveFolder := servePath
	serveFolderAbs, err := filepath.Abs(serveFolder)
	if err != nil {
		log.Fatalf("Invalid serve folder, could not get absolute path: %v", err)
	}

	serveFS := getFileSystems()
	loadTemplates(afero.NewHttpFs(serveFS.templates))

	CreateServer(serveFolderAbs, serveFS, port, portSecondary)
}
