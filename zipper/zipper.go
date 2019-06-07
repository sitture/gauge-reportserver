package zipper

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	OldIndexFilePath = "/index.html"
	NewIndexFilePath = "/report.html"
)

func ZipDir(source, target string) error {

	// check if sourceDir exists
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return err
	}

	source, err := filepath.Abs(source)
	if err != nil {
		return err
	}

	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.Join(strings.TrimPrefix(path, source))

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		// TODO move this out
		// rename file before and after zipping.
		//if path == source + OldIndexFilePath {
		//	err := os.Rename(source + OldIndexFilePath, source + NewIndexFilePath)
		//	if err != nil {
		//		return err
		//	}
		//}

		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)

		return err
	})
	if err != nil {
		return err
	}
	if err = archive.Flush(); err != nil {
		return err
	}
	return nil
}
