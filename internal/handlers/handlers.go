package handlers

import (
	"errors"
	"github.com/Serbroda/jumply/internal/files"
	"github.com/Serbroda/jumply/internal/utils"
	"github.com/Serbroda/jumply/internal/videos"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Handlers struct {
	RootDirs         []string
	VideoFilePattern string
	DefaultPageSize  int
}

func RegisterHandlers(e *echo.Echo, h Handlers, baseUrl string, middlewares ...echo.MiddlewareFunc) {
	e.GET(baseUrl+"/", h.GetIndex, middlewares...)
	e.GET(baseUrl+"/videos/play/:id", h.GetPlay, middlewares...)
	e.GET(baseUrl+"/videos/source/:id", h.GetSource, middlewares...)
	e.GET(baseUrl+"/videos/stream/:id", h.GetStream, middlewares...)
	e.GET(baseUrl+"/reload", h.GetReload, middlewares...)
}

func (h *Handlers) GetIndex(c echo.Context) error {
	if videos.IsEmpty() {
		rg, err := regexp.Compile(h.VideoFilePattern)
		if err != nil {
			panic(err)
		}

		fs := files.ScanAll(h.RootDirs, rg)
		videos.AddAll(utils.MapSlice(fs, func(item files.FileEntry) videos.Video {
			return videos.Video{
				Id:        utils.GenerateID(item.Path),
				FileEntry: item,
			}
		}))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(c.QueryParam("size"))
	if size < 1 {
		size = h.DefaultPageSize
	}

	search := c.QueryParam("search")
	if search != "" && c.QueryParam("page") == "" {
		page = 1
	}

	var items []videos.Video
	if search == "" {
		items = videos.GetAll()
	} else {
		items = videos.Filter(search)
	}

	sortBy := c.QueryParam("sort")
	if sortBy == "" {
		sortBy = "name"
	}
	order := c.QueryParam("order")
	if order == "" {
		order = "asc"
	}
	desc := strings.ToLower(order) == "desc"

	utils.SortSlice(items, func(a, b videos.Video) bool {
		switch sortBy {
		case "size":
			return a.Size < b.Size
		case "modtime":
			return a.ModTime.Before(b.ModTime)
		default:
			return a.Name < b.Name
		}
	}, desc)

	pagination := utils.Paginate(items, page, size)

	if c.Request().Header.Get("HX-Request") != "" {
		// HTMX request: send only partial template
		return c.Render(http.StatusOK, "video_files_partials.html", map[string]any{
			"VideoFiles": pagination,
			"Search":     search,
			"Sort":       sortBy,
			"Order":      order,
		})
	}

	// Normal full page
	return c.Render(http.StatusOK, "index.html", map[string]any{
		"VideoFiles": pagination,
		"Search":     search,
		"Sort":       sortBy,
		"Order":      order,
	})
}

func (h *Handlers) GetPlay(c echo.Context) error {
	id := c.Param("id")
	file, err := videos.GetById(id)
	if err != nil {
		if errors.Is(err, videos.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "video not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.Render(http.StatusOK, "video.html", map[string]any{
		"Video": file,
	})
}

func (h *Handlers) GetSource(c echo.Context) error {
	id := c.Param("id")
	file, err := videos.GetById(id)
	if err != nil {
		if errors.Is(err, videos.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "video not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	p := path.Join(file.Dir, file.Name)
	return c.File(p)
}

func (h *Handlers) GetStream(c echo.Context) error {
	id := c.Param("id")
	file, err := videos.GetById(id)
	if err != nil {
		if errors.Is(err, videos.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "video not found")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	inputPath := path.Join(file.Dir, file.Name)

	if inputPath == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing file parameter")
	}

	// redirect for mp4 files. no transcode necessary
	if strings.HasSuffix(strings.ToLower(file.Name), ".mp4") {
		return c.Redirect(http.StatusTemporaryRedirect, "/videos/source/"+file.Id)
	}

	// check if ffmpeg is available
	_, err = exec.LookPath("ffmpeg")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			"ffmpeg not found: required to stream non-MP4 videos like "+file.Name)
	}

	// ffmpeg transcode
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-f", "mp4",
		"-movflags", "frag_keyframe+empty_moov",
		"-vcodec", "libx264",
		"-acodec", "aac",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentType, "video/mp4")
	c.Response().WriteHeader(http.StatusOK)

	_, err = io.Copy(c.Response(), stdout)
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (h *Handlers) GetReload(c echo.Context) error {
	videos.Clear()
	return c.Redirect(http.StatusPermanentRedirect, "/")
}
