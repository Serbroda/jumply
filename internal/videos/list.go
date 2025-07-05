package videos

import (
	"errors"
	"github.com/Serbroda/jumply/internal/files"
	"github.com/Serbroda/jumply/internal/utils"
	"strings"
)

var (
	ErrNotFound = errors.New("not found")
)

type Video struct {
	Id string
	files.FileEntry
}

var videos = make([]Video, 0)

func Add(video Video) {
	videos = append(videos, video)
}

func AddAll(videos []Video) {
	for _, video := range videos {
		Add(video)
	}
}

func GetById(id string) (Video, error) {
	for _, video := range videos {
		if video.Id == id {
			return video, nil
		}
	}
	return Video{}, ErrNotFound
}

func GetAll() []Video {
	return videos
}

func Filter(search string) []Video {
	return utils.FilterSlice(videos, func(item Video) bool {
		name := strings.ToLower(item.Name)
		s := strings.ToLower(search)

		if strings.Contains(name, s) {
			return true
		}
		if strings.Contains(strings.ReplaceAll(name, ".", " "), s) {
			return true
		}
		return false
	})
}

func Clear() {
	videos = make([]Video, 0)
}

func IsEmpty() bool {
	return len(videos) == 0
}
