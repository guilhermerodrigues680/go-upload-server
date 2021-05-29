package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

// func TestFileHandlerStatusFound(t *testing.T) {
// 	// https://blog.questionable.services/article/testing-http-handlers-go/
// 	req, err := http.NewRequest(http.MethodGet, "/handler", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	s := NewServer("../", false, false, logrus.WithField("test", true))

// 	s.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusFound {
// 		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	if location := rr.HeaderMap.Get("Location"); location != "/handler/" {
// 		t.Fatalf("handler returned wrong header Location: got %v want %v", location, "/handler/")
// 	}
// }

// func TestFileHandlerStatusOK(t *testing.T) {
// 	req, err := http.NewRequest(http.MethodGet, "/handler/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	s := NewServer("../", false, false, logrus.WithField("test", true))

// 	s.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// }

// func TestFileHandlerStream(t *testing.T) {
// 	filepath := "/test/mimetype/yolinux-mime-test.gif"
// 	req, err := http.NewRequest(http.MethodGet, filepath, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	s := NewServer("../", false, false, logrus.WithField("test", true))

// 	s.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	if contentType := rr.Header().Get("Content-Type"); contentType != "image/gif" {
// 		t.Fatalf("handler returned wrong header Location: got %v want %v", contentType, "image/gif")
// 	}

// 	var bufferR bytes.Buffer
// 	name := path.Join(s.staticDirPath, ".") + filepath
// 	buf := make([]byte, 4096)
// 	readFileAndWriteToW(&bufferR, name, buf)

// 	bufferF, err := ioutil.ReadAll(rr.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !bytes.Equal(bufferR.Bytes(), bufferF) {
// 		t.Fatalf("handler returned wrong body: got file size %v want file size %v", bufferR.Len(), len(bufferF))
// 	}
// }

// func TestSpaFileHandlerStream(t *testing.T) {
// 	req, err := http.NewRequest(http.MethodGet, "/random-url-12345", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	// spaMode = true
// 	s := NewServer("../test/spa/dist/", false, true, logrus.WithField("test", true))

// 	s.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	if contentType := rr.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
// 		t.Fatalf("handler returned wrong header Location: got %v want %v", contentType, "text/html; charset=utf-8")
// 	}

// 	var bufferR bytes.Buffer
// 	name := path.Join(s.staticDirPath, "index.html")
// 	buf := make([]byte, 4096)
// 	readFileAndWriteToW(&bufferR, name, buf)

// 	bufferF, err := ioutil.ReadAll(rr.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !bytes.Equal(bufferR.Bytes(), bufferF) {
// 		t.Fatalf("handler returned wrong body: got file size %v want file size %v", bufferR.Len(), len(bufferF))
// 	}
// }

// func TestSpaFileHandlerIndexNotFound(t *testing.T) {
// 	req, err := http.NewRequest(http.MethodGet, "/random-url-12345", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	// spaMode = true
// 	s := NewServer("../test/", false, true, logrus.WithField("test", true))

// 	s.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusNotFound {
// 		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
// 	}
// }

// func TestUploadHandlerStream(t *testing.T) {
// 	s := NewServer("..", false, false, logrus.WithField("test", true))
// 	filepath := "/test/mimetype/yolinux-mime-test.gif"

// 	var b bytes.Buffer
// 	w := multipart.NewWriter(&b)
// 	f, err := os.Open(path.Join(s.staticDirPath, ".") + filepath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer f.Close()

// 	var fw io.Writer
// 	fw, err = w.CreateFormFile("file", path.Base(f.Name()))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if _, err = io.Copy(fw, f); err != nil {
// 		t.Fatal(err)
// 	}

// 	// Don't forget to close the multipart writer.
// 	// If you don't close it, your request will be missing the terminating boundary.
// 	w.Close()

// 	req, err := http.NewRequest(http.MethodPost, "/test/", &b)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// Don't forget to set the content type, this will contain the boundary.
// 	req.Header.Set("Content-Type", w.FormDataContentType())

// 	rr := httptest.NewRecorder()
// 	s.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// }

// func BenchmarkUploadHandlerStream(b *testing.B) {
// 	s := NewServer("..", false, false, logrus.WithField("test", true))
// 	filepath := "/test/mimetype/yolinux-mime-test.gif"

// 	for n := 0; n < b.N; n++ {
// 		var buffer bytes.Buffer
// 		w := multipart.NewWriter(&buffer)
// 		f, err := os.Open(path.Join(s.staticDirPath, ".") + filepath)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 		defer f.Close()

// 		var fw io.Writer
// 		fw, err = w.CreateFormFile("file", path.Base(f.Name()))
// 		if err != nil {
// 			b.Fatal(err)
// 		}

// 		if _, err = io.Copy(fw, f); err != nil {
// 			b.Fatal(err)
// 		}

// 		// Don't forget to close the multipart writer.
// 		// If you don't close it, your request will be missing the terminating boundary.
// 		w.Close()

// 		req, err := http.NewRequest(http.MethodPost, "/test/", &buffer)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 		// Don't forget to set the content type, this will contain the boundary.
// 		req.Header.Set("Content-Type", w.FormDataContentType())

// 		rr := httptest.NewRecorder()
// 		s.ServeHTTP(rr, req)
// 	}

// 	// if status := rr.Code; status != http.StatusOK {
// 	// 	b.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	// }
// }

func BenchmarkFileHandlerStream(b *testing.B) {
	filepath := "/test/mimetype/yolinux-mime-test.gif"
	req, err := http.NewRequest(http.MethodGet, filepath, nil)
	if err != nil {
		b.Fatal(err)
	}

	rr := httptest.NewRecorder()
	s := NewServer("..", false, false, logrus.WithField("test", true))

	for n := 0; n < b.N; n++ {
		s.ServeHTTP(rr, req)
	}
}
