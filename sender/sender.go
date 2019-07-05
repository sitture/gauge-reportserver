package sender

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func SendArchive(url, filePath string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	// Add your image file
	file, err := os.Open(filePath)
	defer func() {
		err := file.Close()
		if err != nil {
			return
		}
	}()
	if err != nil {
		return fmt.Errorf("error opening the file '%s'", filePath)
	}
	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file information")
	}
	formWriter, err := writer.CreateFormFile("file", fileStat.Name())
	if err != nil {
		return fmt.Errorf("error create a form writer")
	}
	if _, err = io.Copy(formWriter, file); err != nil {
		return fmt.Errorf("error copying file '%s' to form writer", file.Name())
	}
	// add auto unzip param to request
	_ = writer.WriteField("unzip", "true")

	// close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	if err = writer.Close(); err != nil {
		return err
	}
	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}
	return
}
