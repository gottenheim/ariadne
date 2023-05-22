package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type Writer struct {
	buffer     bytes.Buffer
	gzipWriter *gzip.Writer
	tarWriter  *tar.Writer
}

func NewWriter() *Writer {
	writer := &Writer{}

	writer.gzipWriter = gzip.NewWriter(&writer.buffer)
	writer.tarWriter = tar.NewWriter(writer.gzipWriter)

	return writer
}

func (w *Writer) AddDir(fs afero.Fs, path string) error {
	return afero.Walk(fs, path, func(file string, fi os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			header.Size = 0
		}

		relPath, err := filepath.Rel(path, file)

		if err != nil {
			return err
		}

		header.Name = relPath

		if err := w.tarWriter.WriteHeader(header); err != nil {
			return err
		}

		err = w.tarWriter.Flush()

		if !fi.IsDir() {
			data, err := fs.Open(file)
			defer data.Close()

			if err != nil {
				return err
			}
			if _, err := io.Copy(w.tarWriter, data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (w *Writer) Buffer() (*bytes.Buffer, error) {
	if err := w.tarWriter.Close(); err != nil {
		return nil, err
	}
	if err := w.gzipWriter.Close(); err != nil {
		return nil, err
	}

	return &w.buffer, nil
}
