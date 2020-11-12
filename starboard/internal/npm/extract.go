package npm

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func createWithNestedDirectories(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

// Adapted from https://stackoverflow.com/questions/57639648/how-to-decompress-tar-gz-file-in-go
func extractTarGz(gzipStream io.Reader, intoFolder string, stripPrefix string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)
	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// We create every file with directories leading up to it, so no need to do anything here.
			// if err := os.Mkdir(path.Join(intoFolder, strings.TrimPrefix(header.Name, stripPrefix)), 0755); err != nil {
			// 	return fmt.Errorf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			// }
		case tar.TypeReg:
			outFile, err := createWithNestedDirectories(path.Join(intoFolder, strings.TrimPrefix(header.Name, stripPrefix)))
			if err != nil {
				return fmt.Errorf("ExtractTarGz: createWithNestedDirectories() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		default:
			return fmt.Errorf(
				"ExtractTarGz: unknown type: %v in %s",
				header.Typeflag,
				header.Name)
		}
	}

	return nil
}
