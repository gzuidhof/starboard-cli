/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/shurcooL/httpgzip"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var indexTemplate *template.Template
var browseTemplate *template.Template
var editorTemplate *template.Template
var fs serveFS

func loadTemplates(fs http.FileSystem) {

	t, err := vfstemplate.ParseGlob(fs, nil, "*.tmpl")
	if err != nil {
		log.Fatalf("Failed to parse templates in vfs: %v", err)
	}
	indexTemplate = t.Lookup("index.html.tmpl")
	browseTemplate = t.Lookup("browse.html.tmpl")
	editorTemplate = t.Lookup("editor.html.tmpl")
}

func Start() {
	port := viper.GetString("port")
	serveFolder := viper.GetString("serve.folder")

	serveFolder, err := filepath.Abs(serveFolder)

	if err != nil {
		log.Fatalf("Invalid serve folder, could not get absolute path: %v", err)
	}

	fs = getFileSystems()
	loadTemplates(fs.templates)

	fileServer := httpgzip.FileServer(fs.static, httpgzip.FileServerOptions{})
	browseHandler := &browseHandler{
		root: http.Dir(serveFolder),
	}

	writeFileSystem := afero.NewBasePathFs(afero.NewOsFs(), serveFolder).(*afero.BasePathFs)
	nbHandler := &notebookHandler{
		root:        http.Dir(serveFolder),
		iframeHost:  "http://localhost:" + port,
		writeFS:     writeFileSystem,
		serveFolder: serveFolder,
	}

	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// /browse/
	http.Handle(defaultBrowseEndpoint, browseHandler)

	// /nb/
	http.Handle(defaultNotebookEndpoint, nbHandler) // Works for both / and /browse/

	//
	http.Handle("/", http.RedirectHandler("/browse", http.StatusFound))

	done := make(chan bool)
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
	log.Printf("Listening on port %v", port)

	<-done
}
