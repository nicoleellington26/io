package io

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// SPA is an http.Handler that serves static files from a http.FileSystem.
// If a static file can not be found, the index file is served.
type SPA struct {
	FileSystem http.FileSystem
	Index      string
	cspHeader  string
}

// EnableCSP generates saves the CSP header to be added to requests of the index file
func (s *SPA) EnableCSP() error {
	var err error
	s.cspHeader, err = s.getCSPHeader()
	return err
}

// ServeHTTP satisfies the http.Handler interface
func (s SPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := s.serveFile(w, r)

	if err == os.ErrNotExist {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.cspHeader != "" {
		contentType := w.Header().Get("Content-Type")
		if strings.Contains(contentType, "text/html") {
			w.Header().Set("Content-Security-Policy", s.cspHeader)
		}
	}

	w.Write(out)
}

func (s SPA) serveFile(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	fileName := getFileName(r, s.Index)

	out, err := ServeFile(w, r, s.FileSystem, fileName, "/")

	// if file does not exist, serve the root index.html entry point
	if os.IsNotExist(err) {
		// If it has an extension, show regular 404
		if filepath.Ext(fileName) != "" {
			return nil, os.ErrNotExist
		}

		out, err = ServeFile(w, r, s.FileSystem, "./"+s.Index, "/")
		// If that does not exist, it's a server error since it always should
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}

		// Successfully served the root index.html file
		return out, nil
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	// Successfully served the requested file
	return out, nil
}

func getFileName(r *http.Request, index string) string {
	fileName := "." + r.URL.Path
	if strings.HasSuffix(fileName, "/") {
		fileName = fileName + index
	}

	return fileName
}
