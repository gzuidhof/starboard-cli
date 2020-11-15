/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"html/template"
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/shurcooL/httpgzip"
	"github.com/spf13/viper"
)

var indexTemplate *template.Template
var browseTemplate *template.Template
var fs serveFS

func loadTemplates(fs http.FileSystem) {

	t, err := vfstemplate.ParseGlob(fs, nil, "*.tmpl")
	if err != nil {
		log.Fatalf("Failed to parse templates in vfs: %v", err)
	}
	indexTemplate = t.Lookup("index.html.tmpl")
	browseTemplate = t.Lookup("browse.html.tmpl")
}

func Start() {
	fs = getFileSystems()
	loadTemplates(fs.templates)

	fileServer := httpgzip.FileServer(fs.static, httpgzip.FileServerOptions{})
	browseHandler := &browseHandler{http.Dir(viper.GetString("serve.folder"))}

	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// /browse/
	http.Handle(defaultBrowseEndpoint, browseHandler) // Works for both / and /browse/
	http.Handle("/", http.RedirectHandler("/browse", http.StatusFound))

	port := viper.GetString("port")
	done := make(chan bool)
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
	log.Printf("Listening on port %v", port)

	<-done
}
