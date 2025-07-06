package main

import (
	"flag"
	"fmt"
	"github.com/Serbroda/jumply"
	"github.com/Serbroda/jumply/internal/files"
	"github.com/Serbroda/jumply/internal/handlers"
	"github.com/Serbroda/jumply/internal/templates"
	"github.com/Serbroda/jumply/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	Version = "dev"

	serverPort  string
	rootDirs    string
	defaultSize int
	videoRegex  string
	customCss   string
)

func main() {
	handleProperties()

	e := echo.New()
	e.Use(middleware.Logger())

	jumply.RegisterStaticFiles(e)

	e.GET("/theme.css", func(c echo.Context) error {
		abs, err := filepath.Abs(customCss)
		if err != nil || !files.FileExists(abs) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.File(abs)
	})

	e.Renderer = templates.NewTemplateRenderer()

	handlers.RegisterHandlers(e, handlers.Handlers{
		RootDirs:         strings.Split(rootDirs, ";"),
		VideoFilePattern: videoRegex,
		DefaultPageSize:  defaultSize,
	}, "")

	printRoutes(e)

	fmt.Printf("Open http://localhost:%s/ in your browser\n", serverPort)
	e.Logger.Fatal(e.Start(":" + serverPort))
}

func handleProperties() {
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Println(Version)
			os.Exit(0)
		}
	}

	utils.StringFlag(&serverPort, utils.GetStringFallback("SERVER_PORT", "8080"), "port", "p", "server port (env: SERVER_PORT)", false)
	utils.StringFlag(&rootDirs, utils.GetStringFallback("ROOT_DIRS", ""), "roots", "r", "root dirs, semicolon-separated (env: ROOT_DIRS)", true)
	utils.IntFlag(&defaultSize, utils.GetInt32Fallback("DEFAULT_PAGE_SIZE", 20), "pagesize", "n", "default page size (env: DEFAULT_PAGE_SIZE)", false)
	utils.StringFlag(&videoRegex, utils.GetStringFallback("VIDEO_FILE_REGEX", `^[^.].*\.(mp4|avi|mkv)$`), "videoregex", "x", "video file regex (env: VIDEO_FILE_REGEX)", false)
	utils.StringFlag(&customCss, utils.GetStringFallback("CUSTOM_CSS_FILE", `./theme.css`), "css", "c", "path to custom CSS (env: CUSTOM_CSS_FILE)", false)

	flag.Parse()
}

func printRoutes(e *echo.Echo) {
	log.Debug("Registered following routes\n\n")
	for _, r := range e.Routes() {
		log.Debugf(" - %v %v\n", r.Method, r.Path)
	}
}
