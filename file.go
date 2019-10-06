package io

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const timeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

// ServeFile from a http.FileSystem
func ServeFile(w http.ResponseWriter, r *http.Request, fileSystem http.FileSystem, fileName string, prefix string) ([]byte, error) {
	index, err := fileSystem.Open(fileName)
	if err != nil {
		return nil, err
	}

	if index == nil {
		return nil, os.ErrNotExist
	}

	stats, err := index.Stat()
	if err != nil {
		return nil, err
	}

	if stats.IsDir() {
		return nil, os.ErrNotExist
	}

	out, err := ioutil.ReadAll(index)
	if err != nil {
		return nil, err
	}

	w.Header().Set("Last-Modified", stats.ModTime().UTC().Format(timeFormat))

	setContentType(w, fileName)
	setCacheControl(w)

	return out, nil
}

func setCacheControl(w http.ResponseWriter) {
	contentType := w.Header().Get("Content-Type")

	// Don't cache HTML files
	if strings.Contains(contentType, "text/html") {
		w.Header().Set("Content-Security-Policy", "script-src 'self'")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=31536000")
}

func setContentType(w http.ResponseWriter, fileName string) {
	extension := filepath.Ext(fileName)
	w.Header().Set("Content-Type", mime.TypeByExtension(extension))
}
