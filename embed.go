package jumply

import (
	"embed"
	"github.com/labstack/echo/v4"
)

var (
	//go:embed internal/templates
	TemplatesDirEmbed embed.FS
	//go:embed internal/static/* internal/static/**/*
	StaticDirEmbed embed.FS

	TemplatesDir = "internal/templates"
	StaticDir    = "internal/static"
)

var (
	templatesDirFS = echo.MustSubFS(TemplatesDirEmbed, TemplatesDir)
	staticDirFS    = echo.MustSubFS(StaticDirEmbed, StaticDir)
)

func RegisterStaticFiles(e *echo.Echo) {
	e.StaticFS("/static", staticDirFS)
}
