/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"

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
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", handleBrowse)

	port := viper.GetString("port")
	done := make(chan bool)
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
	log.Printf("Listening on port %v", port)

	<-done
}

func handleBrowse(w http.ResponseWriter, r *http.Request) {
	loadTemplates(fs.templates) // TODO: Only do this in dev mode
	w.Header().Set("Content-Type", "text/html")

	browsePath := path.Clean(r.URL.Path)
	browsePath = strings.TrimPrefix(browsePath, "/tree")

	var b bytes.Buffer
	err := browseTemplate.Execute(&b, map[string]interface{}{
		"path": browsePath,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = indexTemplate.Execute(w, map[string]interface{}{
		"body": template.HTML(b.String()),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
