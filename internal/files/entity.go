package files

import "time"

type FileEntry struct {
	Path      string
	Name      string
	Size      int64
	Dir       string
	ModTime   time.Time
	Extension string
}
