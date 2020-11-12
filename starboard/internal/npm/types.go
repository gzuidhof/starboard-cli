package npm

import "time"

type NPMRegistryEntry struct {
	ID       string `json:"_id"`
	Rev      string `json:"_rev"`
	Name     string `json:"name"`
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Versions map[string]struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Author      struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Funding struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"funding"`
		License    string `json:"license"`
		Repository struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"repository"`
		Keywords        []string          `json:"keywords"`
		Scripts         map[string]string `json:"scripts"`
		Main            string            `json:"main"`
		DevDependencies map[string]string `json:"devDependencies"`
		Dependencies    map[string]string `json:"dependencies"`
		GitHead         string            `json:"gitHead"`
		ID              string            `json:"_id"`
		NodeVersion     string            `json:"_nodeVersion"`
		NpmVersion      string            `json:"_npmVersion"`
		Dist            struct {
			Integrity    string `json:"integrity"`
			Shasum       string `json:"shasum"`
			Tarball      string `json:"tarball"`
			FileCount    int    `json:"fileCount"`
			UnpackedSize int    `json:"unpackedSize"`
			NpmSignature string `json:"npm-signature"`
		} `json:"dist"`
		Maintainers []struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"maintainers"`
		NpmUser struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"_npmUser"`
		Directories struct {
		} `json:"directories"`
		NpmOperationalInternal struct {
			Host string `json:"host"`
			Tmp  string `json:"tmp"`
		} `json:"_npmOperationalInternal"`
		HasShrinkwrap bool `json:"_hasShrinkwrap"`
	} `json:"versions"`
	Time        map[string]time.Time `json:"time"`
	Maintainers []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"maintainers"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Repository  struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"repository"`
	Author struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
	License        string `json:"license"`
	Readme         string `json:"readme"`
	ReadmeFilename string `json:"readmeFilename"`
}
