/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package nbserver

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
)

var defaultBrowseEndpoint = "/browse/"

// Used for escaping HTML content
var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&#34;",
	"'", "&#39;",
)

// The implementation of the browse handler is based on https://golang.org/src/net/http/fs.go?s=20871:20911#L710

type browseHandler struct {
	root http.FileSystem
}

type BrowseEntry struct {
	Name         string
	URL          string
	IsNotebook   bool
	LastModified string
}

func (f *browseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// loadTemplates(fs.templates) // TODO: Only do this in dev mode
	upath := strings.TrimPrefix(r.URL.Path, defaultBrowseEndpoint)

	log.Print(r)

	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	upath = path.Clean(upath)
	serveFile(w, r, f.root, upath, true)
}

// localRedirect gives a Moved Permanently response.
// It does not convert relative paths to absolute paths like Redirect does.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}

	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}

// name is '/'-separated, not filepath.Separator.
func serveFile(w http.ResponseWriter, r *http.Request, fs http.FileSystem, name string, redirect bool) {
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

	if redirect {
		// redirect to canonical path: / at end of directory url
		// r.URL.Path always begins with /
		url := r.URL.Path

		if d.IsDir() {
			if url[len(url)-1] != '/' {
				localRedirect(w, r, path.Base(url)+"/")
				return
			}
		} else {
			if url[len(url)-1] == '/' {
				localRedirect(w, r, "../"+path.Base(url))
				return
			}
		}
	}

	if d.IsDir() {
		url := r.URL.Path

		// redirect if the directory name doesn't end in a slash
		if url == "" || url[len(url)-1] != '/' {
			localRedirect(w, r, path.Base(url)+"/")
			return
		}

		// if checkIfModifiedSince(r, d.ModTime()) == condFalse {
		// 	writeNotModified(w)
		// 	return
		// }
		// setLastModified(w, d.ModTime())

		dirList(w, r, f)
		return
	}

	// Note(gzuidhof): I replaced serveContent with ServeContent
	// serveContent will check modification time
	// sizeFunc := func() (int64, error) { return d.Size(), nil }
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func dirList(w http.ResponseWriter, r *http.Request, f http.File) {
	dirs, err := f.Readdir(-1)

	if err != nil {
		log.Printf("Error reading directory: %s", err)
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	sort.Slice(dirs, func(i, j int) bool {
		// Folders go above
		if dirs[i].IsDir() && !dirs[j].IsDir() {
			return true
		} else if !dirs[i].IsDir() && dirs[j].IsDir() {
			return false
		}

		return dirs[i].Name() < dirs[j].Name()

	})

	entries := make([]BrowseEntry, len(dirs))
	for i, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		isNotebook := isProbablyNotebookFilename(name)
		URL := url.URL{Path: name}
		if isNotebook {
			URL = url.URL{Path: path.Join(defaultNotebookEndpoint, r.URL.Path, name)}
		}

		entries[i] = BrowseEntry{
			Name:         htmlReplacer.Replace(name),
			URL:          URL.String(),
			IsNotebook:   isNotebook,
			LastModified: d.ModTime().String(),
		}
	}
	crumbs := makeBreadCrumbs(r.URL.Path, true)

	var b bytes.Buffer
	err = browseTemplate.Execute(&b, map[string]interface{}{
		"browseEndpoint":   defaultBrowseEndpoint,
		"notebookEndpoint": defaultNotebookEndpoint,
		"path":             r.URL.Path,
		"breadCrumbs":      crumbs,
		"entries":          entries,
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
