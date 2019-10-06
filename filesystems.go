package io

import (
	"net/http"
	"os"
)

// FileSystems is a wrapper for multiple http.FileSystem
type FileSystems []http.FileSystem

// Open a file from the first http.FileSystem it is found in
func (f FileSystems) Open(name string) (http.File, error) {
	for _, fileSystem := range f {
		file, err := fileSystem.Open(name)
		if os.IsNotExist(err) {
			continue
		}

		return file, err
	}

	return nil, os.ErrNotExist
}
