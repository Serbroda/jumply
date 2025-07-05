package main

import (
	"flag"
	"fmt"
	"github.com/Serbroda/jumply/internal/files"
	"github.com/Serbroda/jumply/internal/templates"
	"github.com/Serbroda/jumply/internal/utils"
	"github.com/Serbroda/jumply/internal/videos"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var (
	Version = "dev"
)

func main() {
	versionFlag := flag.Bool("version", false, "show program version")
	flag.BoolVar(versionFlag, "v", false, "shorthand for --version")
	flag.Parse()

	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	rootDirs := strings.Split(utils.MustGetString("ROOT_DIRS"), ";")
	defaultSize := utils.GetInt32Fallback("DEFAULT_PAGE_SIZE", 20)
	videoFileRegex := utils.GetStringFallback("VIDEO_FILE_REGEX", `^[^.].*\.(mp4|avi|mkv)$`)

	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/static", "internal/static")
	e.Renderer = templates.NewTemplateRenderer()

	e.GET("/", func(c echo.Context) error {
		if videos.IsEmpty() {
			rg, err := regexp.Compile(videoFileRegex)
			if err != nil {
				panic(err)
			}

			fs := files.ScanAll(rootDirs, rg)
			videos.AddAll(utils.MapSlice(fs, func(item files.FileEntry) videos.Video {
				return videos.Video{
					Id:        utils.GenerateID(item.Path),
					FileEntry: item,
				}
			}))
		}

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page < 1 {
			page = 1
		}
		size, err := strconv.Atoi(c.QueryParam("size"))
		if err != nil || size < 1 {
			size = int(defaultSize)
		}

		search := c.QueryParam("search")
		var items []videos.Video
		if search == "" {
			items = videos.GetAll()
		} else {
			items = videos.Filter(search)
		}
		pagination := utils.Paginate(items, page, size)
		return c.Render(http.StatusOK, "index.html", map[string]any{
			"VideoFiles": pagination,
			"Search":     search,
		})
	})

	e.GET("/videos/play", func(c echo.Context) error {
		file := c.QueryParam("file")
		dir := c.QueryParam("dir")
		return c.Render(http.StatusOK, "video.html", map[string]any{
			"VideoName": file,
			"VideoDir":  dir,
		})
	})

	e.GET("/videos/src", func(c echo.Context) error {
		file := c.QueryParam("file")
		dir := c.QueryParam("dir")
		p := path.Join(dir, file)
		return c.File(p)
	})

	e.GET("/videos/stream", func(c echo.Context) error {
		file := c.QueryParam("file")
		dir := c.QueryParam("dir")
		inputPath := path.Join(dir, file)

		if inputPath == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing file parameter")
		}

		// redirect for mp4 files. no transcode necessary
		if strings.HasSuffix(strings.ToLower(file), ".mp4") {
			return c.Redirect(http.StatusTemporaryRedirect, "/videos/src?file="+file+"&dir="+dir)
		}

		// check if ffmpeg is available
		_, err := exec.LookPath("ffmpeg")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError,
				"ffmpeg not found: required to stream non-MP4 videos like "+file)
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
	})

	e.GET("/reload", func(c echo.Context) error {
		videos.Clear()
		return c.Redirect(http.StatusPermanentRedirect, "/")
	})

	printRoutes(e)
	e.Logger.Fatal(e.Start(":" + utils.GetStringFallback("SERVER_PORT", "8080")))
}

func printRoutes(e *echo.Echo) {
	log.Debug("Registered following routes\n\n")
	for _, r := range e.Routes() {
		log.Debugf(" - %v %v\n", r.Method, r.Path)
	}
}
