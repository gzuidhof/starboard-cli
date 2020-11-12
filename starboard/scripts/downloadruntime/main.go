package main

import (
	"log"
	"os"

	"github.com/gzuidhof/starboard-cli/starboard/internal/npm"
)

const defaultRuntimePackageName = "starboard-notebook"

func main() {
	if len(os.Args) < 3 {
		log.Print("Not enough arguments, supply 2 arguments: the version and output folder")
		os.Exit(1)
	}
	version := os.Args[1]
	outFolder := os.Args[2]

	id, err := npm.DownloadPackageIntoFolder(defaultRuntimePackageName, version, outFolder)

	if err != nil {
		log.Fatalf("Failed to fetch %s: %v", defaultRuntimePackageName, err)
	}
	log.Printf("Downloaded %s into %s", id, outFolder)
}
