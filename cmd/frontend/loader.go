package main

import (
	"embed"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/CloudyKit/jet/v6"
)

type embedFileSystemLoader struct {
	dir string
	fs  embed.FS
}

// NewLoader returns an initialized loader serving the passed embed.FS.
func NewLoader(dirPath string, fs embed.FS) jet.Loader {
	return &embedFileSystemLoader{
		dir: filepath.FromSlash(dirPath),
		fs:  fs,
	}
}

// Open implements Loader.Open() on top of an embed.FS.
func (l *embedFileSystemLoader) Open(name string) (io.ReadCloser, error) {
	return l.fs.Open(filepath.Join(l.dir, filepath.FromSlash(name)))
}

// Exists implements Loader.Exists() on top of an embed.FS by trying to open the file.
func (l *embedFileSystemLoader) Exists(name string) bool {
	name = filepath.Join(l.dir, filepath.FromSlash(name))
	stat, err := fs.Stat(l.fs, name)
	if err == nil && !stat.IsDir() {
		return true
	}
	return false
}
