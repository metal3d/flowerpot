package petalsserver

import "os"

func EmptyCache() error {
	cacheDir := os.ExpandEnv(cacheDir)
	return os.RemoveAll(cacheDir)
}
