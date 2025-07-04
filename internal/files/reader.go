package files

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Scan(root string, pattern *regexp.Regexp) ([]FileEntry, error) {
	var videos []FileEntry

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if !pattern.MatchString(strings.ToLower(info.Name())) {
			return nil
		}

		relativePath, err := filepath.Rel(root, path)

		if err != nil {
			return nil // skip file on ffprobe error
		}

		videos = append(videos, FileEntry{
			Path: relativePath,
			Name: info.Name(),
			Size: info.Size(),
			Dir:  root,
		})
		return nil
	})

	return videos, err
}
