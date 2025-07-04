package templates

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	_ "path/filepath"
	"time"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"now": func() string {
			now := time.Now()
			return now.Format("20060102150405") + fmt.Sprintf("%04d", now.Nanosecond()/1e5)
		},
	}).ParseGlob("internal/templates/*.html"))
	return &TemplateRenderer{templates: tmpl}
}

func (r *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}
