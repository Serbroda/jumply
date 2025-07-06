package templates

import (
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/Serbroda/jumply"
	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer lädt Templates aus dem eingebetteten FS.
func NewTemplateRenderer() *TemplateRenderer {
	tmpl := template.Must(
		template.New("").Funcs(template.FuncMap{
			"add": func(a, b int) int { return a + b },
			"now": func() string {
				n := time.Now()
				return n.Format("20060102150405") + fmt.Sprintf("%04d", n.Nanosecond()/1e5)
			},
		}).ParseFS(jumply.TemplatesDirEmbed, "internal/templates/*.html"),
	)
	return &TemplateRenderer{templates: tmpl}
}

func (r *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}
