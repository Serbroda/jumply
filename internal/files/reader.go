package files

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return false
	}
}

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

		if err != nil {
			return nil
		}

		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return nil
		}

		videos = append(videos, FileEntry{
			Name:      info.Name(),
			Path:      absolutePath,
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			Extension: filepath.Ext(path),
		})
		return nil
	})

	return videos, err
}

func ScanAll(rootDirs []string, rg *regexp.Regexp) []FileEntry {
	var videos []FileEntry
	for _, rootDir := range rootDirs {
		vs, err := Scan(rootDir, rg)
		if err == nil {
			for _, v := range vs {
				videos = append(videos, v)
			}
		}
	}
	return videos
}
