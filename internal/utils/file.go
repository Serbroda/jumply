package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type DirEntry struct {
	Name                 string    `json:"name"`
	NameWithoutExtension string    `json:"nameWithoutExtension"`
	Path                 string    `json:"path"`
	Dir                  string    `json:"dir"`
	Extension            string    `json:"ext"`
	IsDir                bool      `json:"isDir"`
	ModTime              time.Time `json:"modTime"`
	Size                 int64     `json:"size"`
	Hash                 uint32    `json:"hash"`
}

type ReadDirOptions struct {
	Recursive   bool
	ExcludeDirs bool
	NamePattern []string
}

func FilenameWithoutExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return false
	}
}

func FileInfo(path string) (DirEntry, error) {
	file, err := os.Stat(path)
	if err != nil {
		return DirEntry{}, err
	}

	return DirEntry{
		Name:                 file.Name(),
		NameWithoutExtension: FilenameWithoutExtension(file.Name()),
		Path:                 path,
		Dir:                  filepath.Dir(path),
		IsDir:                file.IsDir(),
		ModTime:              file.ModTime(),
		Size:                 file.Size(),
		Extension:            filepath.Ext(path),
		Hash:                 Hash(path),
	}, nil
}

func ReadDir(dir string) ([]DirEntry, error) {
	return ReadDirOpt(dir, ReadDirOptions{
		Recursive:   false,
		ExcludeDirs: false,
		NamePattern: []string{},
	})
}

func ReadDirOpt(dir string, opt ReadDirOptions) ([]DirEntry, error) {
	entries := make([]DirEntry, 0)
	return walkDir(dir, &entries, opt)
}

func walkDir(dir string, entries *[]DirEntry, opt ReadDirOptions) ([]DirEntry, error) {
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return []DirEntry{}, err
	}

	for _, itm := range items {
		absolute := filepath.Join(dir, itm.Name())
		isDir := itm.IsDir()

		if opt.Recursive && isDir {
			walkDir(absolute, entries, opt)
		}

		if opt.ExcludeDirs && isDir {
			continue
		}

		file := DirEntry{
			Name:      itm.Name(),
			Path:      absolute,
			Dir:       filepath.Dir(absolute),
			IsDir:     isDir,
			ModTime:   itm.ModTime(),
			Size:      itm.Size(),
			Extension: filepath.Ext(absolute),
			Hash:      Hash(absolute),
		}

		if isDir || len(opt.NamePattern) < 1 {
			*entries = append(*entries, file)
		} else {
			for _, p := range opt.NamePattern {
				match, err := regexp.MatchString(p, file.Name)

				if err != nil {
					fmt.Printf("warn: failed to check pattern: %v\n", err.Error())
				} else if match {
					*entries = append(*entries, file)
				}
			}
		}
	}
	return *entries, nil
}
