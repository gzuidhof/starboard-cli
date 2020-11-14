/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gzuidhof/starboard-cli/starboard/assets/web_static"
	"github.com/gzuidhof/starboard-cli/starboard/assets/web_template"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/shurcooL/httpgzip"
	"github.com/spf13/viper"
)

var indexTemplate *template.Template

func initTemplates() {
	t, err := vfstemplate.ParseGlob(web_template.TemplateAssets, nil, "*.tmpl")
	if err != nil {
		log.Fatalf("Failed to parse templates in vfs: %v", err)
	}
	indexTemplate = t.Lookup("index.html.tmpl")
}

func Start() {
	initTemplates()

	fs := httpgzip.FileServer(web_static.StaticAssets, httpgzip.FileServerOptions{})
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleHome)

	port := viper.GetString("port")
	done := make(chan bool)
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()
	log.Printf("Listening on port %v", port)

	<-done
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := indexTemplate.Execute(w, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
