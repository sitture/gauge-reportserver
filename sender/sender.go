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
	if err != nil {
		return
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}

	defer file.Close()

	fw, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return err
	}

	// Add the other fields

	//if fw, err = w.CreateFormField("key"); err != nil {
	//	return
	//}
	//if _, err = fw.Write([]byte("KEY")); err != nil {
	//	return
	//}

	_ = writer.WriteField("unzip", "true")

	// Don't forget to close the multipart writer.
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
		err = fmt.Errorf("bad status: %s", res.Status)
	}

	RemoveArchive(filePath)

	return
}

func RemoveArchive(filePath string) {
	if err := os.Remove(filePath); err != nil {
		return
	}
}
