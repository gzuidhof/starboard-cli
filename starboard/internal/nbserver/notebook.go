/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/spf13/afero"
)

var defaultNotebookEndpoint = "/nb/"

type notebookHandler struct {
	iframeHost  string
	root        http.FileSystem
	writeFS     *afero.BasePathFs
	serveFolder string
}

func (h *notebookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	loadTemplates(fs.templates) // TODO: Only do this in dev mode
	if r.Method == http.MethodGet {
		upath := strings.TrimPrefix(r.URL.Path, defaultNotebookEndpoint)
		r.URL.Path = upath

		// if !strings.HasPrefix(upath, "/") {
		// 	upath = "/" + upath
		// 	r.URL.Path = upath
		// }

		upath = path.Clean(upath)
		h.serveNotebook(w, r, h.root, upath)
	} else if r.Method == http.MethodPut {
		upath := strings.TrimPrefix(r.URL.Path, defaultNotebookEndpoint)
		upath = path.Clean(upath)

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read body, bad request.", http.StatusBadRequest)
			return
		}

		log.Printf("Writing to %s (%v bytes)", upath, len(content))
		err = afero.WriteFile(h.writeFS, upath, content, 0644)
		if err != nil {
			log.Printf("Failed to save notebook to disk: %v", err)
			http.Error(w, "Could not save notebook", http.StatusInternalServerError)
			return
		}
	}
}

// name is '/'-separated, not filepath.Separator.
func (h *notebookHandler) serveNotebook(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string) {
	f, err := fs.Open(name)

	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	defer f.Close()
	d, err := f.Stat()

	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	if d.IsDir() {
		url := r.URL.Path
		// redirect if the directory name doesn't end in a slash
		if url == "" || url[len(url)-1] != '/' {
			localRedirect(w, r, path.Base(url)+"/")
			return
		}

		localRedirect(w, r, path.Join(defaultBrowseEndpoint, url))
		return
	}

	fileContent, err := ioutil.ReadAll(f)

	if err != nil {
		log.Print("Failed to read file")
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	crumbs := makeBreadCrumbs(r.URL.Path, false)

	var b bytes.Buffer
	err = editorTemplate.Execute(&b, map[string]interface{}{
		"browseEndpoint":   defaultBrowseEndpoint,
		"notebookEndpoint": defaultNotebookEndpoint,
		"path":             r.URL.Path,
		"breadCrumbs":      crumbs,
		"iframeHost":       h.iframeHost,
		"notebookContent":  string(fileContent),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = indexTemplate.Execute(w, map[string]interface{}{
		"body": template.HTML(b.String()),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
