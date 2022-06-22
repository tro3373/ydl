package middleware

import "net/http"

type statikFileSystem struct {
	fs http.FileSystem
}

func (b *statikFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *statikFileSystem) Exists(prefix string, filepath string) bool {
	if _, err := b.fs.Open(filepath); err != nil {
		return false
	}
	return true
}

func StatikFileSystem(fs http.FileSystem) *statikFileSystem {
	return &statikFileSystem{
		fs,
	}
}
