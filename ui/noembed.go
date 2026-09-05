//go:build !embedui

package ui

import (
	"io/fs"
	"os"
)

// DistFS falls back to serving ui/dist from the local filesystem when running in development
// without the `embedui` build tag.
func DistFS() (fs.FS, bool) {
	if _, err := os.Stat("ui/dist"); err == nil {
		return os.DirFS("ui/dist"), true
	}
	return nil, false
}
