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

// Deletes the src and test folder in the output, saves some KB in executable size.
func deleteUselessFiles(fromFolder string) {
	// Delete src folder
	err := os.RemoveAll(path.Join(fromFolder, "dist/src"))
	if err != nil {
		log.Fatalf("Failed to delete: %v", err)
	}
	err = os.RemoveAll(path.Join(fromFolder, "dist/test"))
	if err != nil {
		log.Fatalf("Failed to delete: %v", err)
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
	if len(os.Args) < 3 {
		log.Print("Not enough arguments, supply 2 arguments: the version and output folder")
		os.Exit(1)
	}
	version := os.Args[1]
	outFolder := os.Args[2]

	if version != "latest" && dirExists(path.Join(outFolder, defaultRuntimePackageName+"@"+version)) {
		log.Printf("Skipping NPM fetch of %s as version %v already seems to be vendored already", defaultRuntimePackageName, version)
		os.Exit(0)
	}

	// id has the form <packagename>@<version>
	id, err := npm.DownloadPackageIntoFolder(defaultRuntimePackageName, version, outFolder)

	packageFolder := path.Join(outFolder, id)
	deleteUselessFiles(packageFolder)

	if err != nil {
		log.Fatalf("Failed to fetch %s: %v", defaultRuntimePackageName, err)
	}
	log.Printf("Downloaded %s into %s", id, outFolder)
}
