package main

import (
	"flag"
	"fmt"
	"github.com/Serbroda/jumply/internal/handlers"
	"github.com/Serbroda/jumply/internal/templates"
	"github.com/Serbroda/jumply/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"os"
	"strings"
)

var (
	Version = "dev"

	serverPort     = utils.GetStringFallback("SERVER_PORT", "8080")
	rootDirs       = utils.MustGetString("ROOT_DIRS")
	defaultSize    = utils.GetInt32Fallback("DEFAULT_PAGE_SIZE", 20)
	videoFileRegex = utils.GetStringFallback("VIDEO_FILE_REGEX", `^[^.].*\.(mp4|avi|mkv)$`)
)

func init() {
	versionFlag := flag.Bool("version", false, "show program version")
	flag.BoolVar(versionFlag, "v", false, "shorthand for --version")
	flag.Parse()

	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Static("/static", "internal/static")
	e.Renderer = templates.NewTemplateRenderer()

	handlers.RegisterHandlers(e, handlers.Handlers{
		RootDirs:         strings.Split(rootDirs, ";"),
		VideoFilePattern: videoFileRegex,
		DefaultPageSize:  defaultSize,
	}, "")

	printRoutes(e)
	e.Logger.Fatal(e.Start(":" + serverPort))
}

func printRoutes(e *echo.Echo) {
	log.Debug("Registered following routes\n\n")
	for _, r := range e.Routes() {
		log.Debugf(" - %v %v\n", r.Method, r.Path)
	}
}
