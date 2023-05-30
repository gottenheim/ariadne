package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

func Uncompress(src io.Reader, fs afero.Fs, targetDir string) error {
	zr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}

	tr := tar.NewReader(zr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		// validate name against path traversal
		if !validRelPath(header.Name) {
			return fmt.Errorf("tar contained invalid name error %q\n", header.Name)
		}

		// add dst + re-format slashes according to system
		target := filepath.Join(targetDir, header.Name)
		// if no join is needed, replace with ToSlash:
		// target = filepath.ToSlash(header.Name)

		// check the type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it (with 0755 permission)
		case tar.TypeDir:
			if _, err := fs.Stat(target); err != nil {
				if err := fs.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// if it's a file create it (with same permission)
		case tar.TypeReg:
			fileToWrite, err := fs.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return err
			}
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			err = fileToWrite.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetDirectories(src io.Reader) ([]string, error) {
	zr, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(zr)

	var archiveDirs []string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, err
		}

		if header.Typeflag == tar.TypeDir {
			archiveDirs = append(archiveDirs, header.Name)
		}
	}

	return archiveDirs, nil
}

func GetFiles(src io.Reader) (map[string][]byte, error) {
	zr, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(zr)

	archiveFiles := map[string][]byte{}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, err
		}

		if header.Typeflag == tar.TypeReg {
			buf := new(bytes.Buffer)
			_, err := io.Copy(buf, tr)
			if err != nil {
				return nil, err
			}
			archiveFiles[header.Name] = buf.Bytes()
		}
	}

	return archiveFiles, nil
}

func GetContentHash(src io.Reader) (string, error) {
	zr, err := gzip.NewReader(src)
	if err != nil {
		return "", err
	}

	tr := tar.NewReader(zr)

	buffer := &bytes.Buffer{}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return "", err
		}

		if header.Typeflag == tar.TypeReg {
			_, err := io.Copy(buffer, tr)
			if err != nil {
				return "", err
			}
		}
	}

	hash := md5.Sum(buffer.Bytes())

	return fmt.Sprintf("%x", hash), nil
}

// check for path traversal and correct forward slashes
func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}
