package zipper

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	zipFile, err := os.Create(target)
	defer func() {
		err := zipFile.Close()
		if err != nil {
			return
		}
	}()
	if err != nil {
		return err
	}

	archive := zip.NewWriter(zipFile)
	defer func() {
		err := archive.Close()
		if err != nil {
			return
		}
	}()

	ignoreFiles := []string{
		"html-report",
		".DS_Store",
		"report.html",
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if contains(ignoreFiles, info.Name()) {
			return nil
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

		file, err := os.Open(path)
		defer func() {
			err := file.Close()
			if err != nil {
				return
			}
		}()
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

func contains(array []string, str string) bool {
	for _, a := range array {
		if a == str {
			return true
		}
	}
	return false
}
