package handler

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Server struct {
	r                          *httprouter.Router
	logger                     *logrus.Entry
	staticDirPath              string
	keepOriginalUploadFileName bool
	spaMode                    bool
}

func NewServer(staticDirPath string, keepOriginalUploadFileName bool, spaMode bool, logger *logrus.Entry) *Server {
	router := httprouter.New()
	s := Server{
		r:                          router,
		logger:                     logger,
		staticDirPath:              staticDirPath,
		keepOriginalUploadFileName: keepOriginalUploadFileName,
		spaMode:                    spaMode,
	}

	if s.spaMode {
		router.GET("/*filepath", s.spaFileHandler)
	} else {
		router.GET("/*filepath", s.fileHandler)
		router.POST("/*dirpath", s.uploadHandler)
	}

	return &s
}

func (f *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mw := NewLoggingInterceptorOnServer(f.r, f.logger.WithField("server", "interceptor-on-server"))
	mw.ServeHTTP(w, r)
}

func (s *Server) fileHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fileUrlPath := p.ByName("filepath")
	filePath := path.Join(s.staticDirPath, ".", fileUrlPath)
	s.logger.Trace(filePath)

	fileinfo, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	switch mode := fileinfo.Mode(); {
	case mode.IsDir():
		// isDir but not HasSuffix '/'
		if !strings.HasSuffix(fileUrlPath, "/") {
			http.Redirect(w, r, fileUrlPath+"/", http.StatusFound)
			return
		}

		err := sendDirFileListToClient(w, filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case mode.IsRegular():
		err := sendFileToClient(w, filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("Error: Unrecognized mode %s", mode), http.StatusInternalServerError)
	}
}

// spaFileHandler is a handler for Single Page Aplications (SPA).
// Single Page Aplications (SPA), são aplicações que possuem uma única página o 'index.html'
// Essas aplicações podem usar a History API do HTML5 para simular uma navegação entre páginas.
// Porém essa abordagem tem um problema que é quando o usuário faz o refresh na página, pois como
// a rota no histórico foi programada, ela não existirá no servidor. Assim para atender essas
// aplicações é necessário que o servidor não envie um Status 404 Not Found nessas situações e
// sim o proprio 'index.html', pois a SPA se encaregará de renderizar a página correta ou exibir
// um erro.
// Ex: https://router.vuejs.org/guide/essentials/history-mode.html
func (s *Server) spaFileHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	indexPath := path.Join(s.staticDirPath, "./index.html")

	fileUrlPath := p.ByName("filepath")
	filePath := path.Join(s.staticDirPath, ".", fileUrlPath)

	// root requests receive the 'index.html' file
	if fileUrlPath == "/" || fileUrlPath == "" {
		filePath = indexPath
	}

	s.logger.Trace(filePath)

	err := sendFileToClient(w, filePath)
	if err == nil {
		// FIXME - cliente broken pipe
		// OK! file successfully sent to the client
		return
	}

	if !errors.Is(err, os.ErrNotExist) {
		// unknown error returns an internal server error
		s.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// could not find the file path, fallback to 'index.html'
	s.logger.Infof("%s Not Found. Responding to the request with the index.html", filePath)
	err = sendFileToClient(w, indexPath)
	if err == nil {
		// OK! file successfully sent to the client
		return
	}

	if !errors.Is(err, os.ErrNotExist) {
		// unknown error returns an internal server error
		s.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 'index.html' not found returns a 404 status
	http.Error(w, err.Error(), http.StatusNotFound)
}

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	dirUrlPath := p.ByName("dirpath")
	dirPath := path.Join(s.staticDirPath, ".", path.Dir(dirUrlPath))
	s.logger.Trace(dirPath)

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		s.logger.Errorf("Parse Media Type error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	s.logger.Trace(mediaType, params)

	boundary := params["boundary"]
	reader := multipart.NewReader(r.Body, boundary)
	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				// Done reading body
				break
			}
			s.logger.Errorf("Multipart Reader NextPart error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// only accept fieldname == "file", otherwise, return a validation err
		if part.FormName() != "file" {
			s.logger.Errorf("Field Name != 'file'. Got %s", part.FormName())
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Field Name != 'file'. Got %s", part.FormName())
			return
		}

		contentType := part.Header.Get("Content-Type") // FIXME Not used
		fname := part.FileName()
		s.logger.Infof("multipart/form-data Content-Type: %s, Filename: %s", contentType, fname)

		// bufferedWriter := bufio.NewWriter(w)
		buf := make([]byte, 4096) // make a buffer to keep chunks that are read
		fileSent, err := readerToFile(part, dirPath, fname, s.keepOriginalUploadFileName, buf)
		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) {
				s.logger.Errorf("Reader To File error, Client closed the connection: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			} else {
				s.logger.Errorf("Reader To File error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
			return
		}
		s.logger.Infof("File sent: %s", fileSent)
	}
}

// helpers

func getContentType(f *os.File) (string, error) {
	ctype := mime.TypeByExtension(filepath.Ext(f.Name()))
	if ctype == "" {
		const sniffLen = 512
		var buf [sniffLen]byte
		n, _ := io.ReadFull(f, buf[:]) // read a chunk to decide between utf-8 text and binary
		ctype = http.DetectContentType(buf[:n])
		_, err := f.Seek(0, io.SeekStart) // rewind to output whole file
		if err != nil {
			// seeker can't seek
			return "", err
		}
	}

	return ctype, nil
}

func formatBytes(b int64) string {
	const unit = 1024 // ByteCountIEC
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func sendFileToClient(w http.ResponseWriter, filepath string) error {
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	if !fileinfo.Mode().IsRegular() {
		return ErrFileIsNotRegular
	}

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("os.Open(filepath) error: %w", err)
	}
	defer file.Close()

	fmt.Println(file.Name())
	ctype, err := getContentType(file)
	if err != nil {
		return fmt.Errorf("Get Content-Type error: %w", err)
	}

	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Content-Length", strconv.FormatInt(fileinfo.Size(), 10))
	w.WriteHeader(http.StatusOK)

	bufferedReader := bufio.NewReader(file)

	n, err := io.Copy(w, bufferedReader)
	if err != nil {
		fmt.Println(err.Error()) // FIXME
		return fmt.Errorf("io.Copy Error: %w", err)
	}

	fmt.Printf("n: %d, fileinfo.Size: %d\n", n, fileinfo.Size()) // FIXME

	return nil
}

func sendDirFileListToClient(w http.ResponseWriter, dirpath string) error {
	fileinfo, err := os.Stat(dirpath)
	if err != nil {
		return err
	}

	if !fileinfo.Mode().IsDir() {
		return ErrFileIsNotDir
	}

	dirfileList, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}

	// Sort by name
	sort.Slice(dirfileList, func(i, j int) bool {
		return strings.ToLower(dirfileList[i].Name()) < strings.ToLower(dirfileList[j].Name())
	})

	t, err := template.New("files").Funcs(template.FuncMap{
		"formatBytes": formatBytes,
	}).Parse(TemplateListFiles)
	if err != nil {
		return fmt.Errorf("%w %s", ErrCreateTemplate, err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = t.Execute(w, dirfileList)
	if err != nil {
		return fmt.Errorf("%w %s", ErrExecuteTemplate, err)
	}

	return nil
}

func readerToFile(r io.Reader, dir string, fname string, keepOriginalFileName bool, buf []byte) (string, error) {
	// FIXME file permissions originais

	ext := path.Ext(fname)
	name := fname[0 : len(fname)-len(ext)]
	tempFilePattern := name + "-*" + ext
	tempFile, err := ioutil.TempFile(dir, tempFilePattern)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			os.Remove(tempFile.Name())
			return "", err
		}

		if n == 0 {
			break
		}

		// write a chunk
		if _, err := tempFile.Write(buf[:n]); err != nil {
			os.Remove(tempFile.Name())
			return "", err
		}
	}

	if keepOriginalFileName {
		finalFileName := path.Join(dir, fname)
		err = os.Rename(tempFile.Name(), path.Join(dir, fname))
		if err != nil {
			return "", err
		}
		return finalFileName, nil
	}

	return tempFile.Name(), nil
}
