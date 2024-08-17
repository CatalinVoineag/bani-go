package main

import (
	"html/template"
	"io"
//	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
  templates *template.Template
}

func (t *Templates) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
  return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
  return &Templates {
    templates: template.Must(template.ParseGlob("views/*.html")),
  }
}
type Page struct {
  Data string 
}

func newPage() Page {
  return Page{
    Data: "hi",
  }
}

func main() {
  e := echo.New()
  e.Use(middleware.Logger())

  page := newPage()
  e.Renderer = newTemplate()

  e.GET("/", func(c echo.Context) error {
    return c.Render(200, "index", page)
  })

  e.Logger.Fatal(e.Start(":3000"))
}
