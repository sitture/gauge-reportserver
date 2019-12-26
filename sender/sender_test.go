package sender

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSenderSendArchiveWithCorrectRequest(t *testing.T) {
	expectedMethod := "POST"
	expectedContentTypePart := "multipart/form-data"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != expectedMethod {
			t.Fatalf("expected http method %s but was %s", expectedMethod, r.Method)
		}

		contentType := r.Header.Get("Content-Type")

		if !strings.Contains(contentType, expectedContentTypePart) {
			t.Fatalf("expected Content-Type to contain %s. actual value: %s", expectedContentTypePart, contentType)
		}

		fmt.Fprintln(w, "ok")
	}))

	defer ts.Close()

	err := SendArchive(ts.URL, "testdata/test.txt")

	if err != nil {
		t.Fatal(err)
	}

}


func TestSenderSendArchiveHttpServerError(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))

	defer ts.Close()

	err := SendArchive(ts.URL, "testdata/test.txt")

	if err == nil {
		t.Fatal("expected error on http server status code 500")
	}

}
