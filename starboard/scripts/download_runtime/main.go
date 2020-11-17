/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"os"
	"path"

	"github.com/gzuidhof/starboard-cli/starboard/internal/npm"
)

const defaultRuntimePackageName = "starboard-notebook"
const defaultIframeResizerPackageName = "iframe-resizer"

// Deletes the src and test folder in the output, saves some KB in executable size.
func deleteUselessFiles(fromFolder string, toRemove []string) {
	for _, folder := range toRemove {
		err := os.RemoveAll(path.Join(fromFolder, folder))
		if err != nil {
			log.Fatalf("Failed to delete: %v", err)
		}
	}
}

func dirExists(dirpath string) bool {
	info, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func main() {
	if len(os.Args) < 4 {
		log.Print("Not enough arguments, supply 3 arguments: the package name, version and output folder")
		os.Exit(1)
	}
	packageName := os.Args[1]
	version := os.Args[2]
	outFolder := os.Args[3]

	if version != "latest" && dirExists(path.Join(outFolder, packageName+"@"+version)) {
		log.Printf("Skipping NPM fetch of %s as version %v already seems to be vendored already", packageName, version)
		os.Exit(0)
	}

	// id has the form <packagename>@<version>
	id, err := npm.DownloadPackageIntoFolder(packageName, version, outFolder)
	packageFolder := path.Join(outFolder, id)

	if packageName == defaultRuntimePackageName {
		deleteUselessFiles(packageFolder, []string{"dist/src", "dist/test"})
	} else if packageName == defaultIframeResizerPackageName {
		// Possible improvement: create a "keep-these-files" function.. we really only need one file
		deleteUselessFiles(packageFolder, []string{
			".github",
			"CHANGELOG.md", // Some large markdown files we really don't have to statically bundle
			"CONTRIBUTING.md",
			"README.md",
			"js/iframeResizer.js", // We use the minified version instead
			"js/iframeResizer.contentWindow.js",
			"js/iframeResizer.contentWindow.map",
			"js/iframeResizer.contentWindow.min.js",
			"js/index.js",
		})
	}

	if err != nil {
		log.Fatalf("Failed to fetch %s: %v", packageName, err)
	}
	log.Printf("Downloaded %s into %s", id, outFolder)
}
