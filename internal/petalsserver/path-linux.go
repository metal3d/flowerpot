//go:build linux

package petalsserver

import (
	"os"
	"path/filepath"
)

const (
	shareDir   = "${HOME}/.local/share"
	installDir = "petals-server"
	cacheDir   = "${HOME}/.cache/petals"
)

func pipPath() string {
	return os.ExpandEnv(filepath.Join(shareDir, installDir, "bin", "pip"))
}

func pythonPath() string {
	return os.ExpandEnv(filepath.Join(shareDir, installDir, "bin", "python"))
}
