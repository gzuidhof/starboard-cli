/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package npm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

const defaultNPMRegistryEndpoint = "https://registry.npmjs.org/"

func downloadAndUnpackTarball(tarballURL string, intoFolder string) error {
	log.Printf("Downloading tarball from %s", tarballURL)
	resp, err := http.Get(tarballURL)
	if err != nil {
		return fmt.Errorf("failed to fetch tarball from %s: %v", tarballURL, err)
	}
	defer resp.Body.Close()

	return extractTarGz(resp.Body, intoFolder, "package/")
}

func GetRegistryEntry(packageName string) (NPMRegistryEntry, error) {
	url := defaultNPMRegistryEndpoint + packageName
	resp, err := http.Get(url)
	if err != nil {
		return NPMRegistryEntry{}, fmt.Errorf("failed to fetch %s: %v", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var sb NPMRegistryEntry

	err = json.Unmarshal(body, &sb)
	if err != nil {
		return NPMRegistryEntry{}, fmt.Errorf("failed to unmarshal NPM response: %v", err)
	}
	return sb, nil
}

/**
* DownloadPackageIntoFolder downloads given NPM package into the outfolder, it puts it in a subfolder
* called <packageName>@<version>, this function returns that folder name.
 */
func DownloadPackageIntoFolder(packageName string, version string, outFolder string) (string, error) {

	sb, err := GetRegistryEntry(packageName)
	if err != nil {
		return "", fmt.Errorf("failed to get registry entry: %v", err)
	}

	versionToGet := version
	if version == "latest" {
		versionToGet = sb.DistTags.Latest
	}

	packageID := packageName + "@" + versionToGet

	outFolder = path.Join(outFolder, packageID)

	tarballURL := sb.Versions[versionToGet].Dist.Tarball
	err = downloadAndUnpackTarball(tarballURL, outFolder)
	if err != nil {
		return "", fmt.Errorf("failed to download and handle tarball: %v", err)
	}

	return packageID, nil
}
