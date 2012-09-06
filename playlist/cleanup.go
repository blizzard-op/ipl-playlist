package playlist

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func Cleanup() error {
	tmpFilename := "ipl-playlist*"
	var pattern string
	if(runtime.GOOS == "windows"){
		a := []string{os.TempDir(), tmpFilename}
		pattern = strings.Join(a, "\\")
	} else {
		pattern = path.Join(os.TempDir(), tmpFilename)
	}
	filepaths, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, p := range filepaths {
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
		err = os.Remove(f.Name())
		if err != nil {
			return err
		}
	}
	return nil
}