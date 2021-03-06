package zipper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ZipDir(source, target string) error {
	// check if sourceDir exists
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return fmt.Errorf("source directory '%s' does not exist", source)
	}
	// get absolute path of source directory
	source, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("could not get absolute path for the source directory '%s'", source)
	}
	// create a zip file at target
	zipFile, err := os.Create(target)
	defer zipFile.Close()
	if err != nil {
		return fmt.Errorf("could not create the target archive at '%s'", target)
	}
	// create a new writer for writing into zip
	archive := zip.NewWriter(zipFile)
	defer archive.Close()
	// list of files to ignore when adding to zip
	ignoreFiles := []string{
		"html-report",
		".DS_Store",
	}
	if err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// check file isn't in the ignored list
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
		// skip if it's a directory
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
	}); err != nil {
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
