package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"meal-choices/db"
	"meal-choices/routes"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	templates, err := findAndParseTemplates("views", nil)

	if err != nil {
		panic(err)
	}

	return &Templates{
		templates: templates,
	}
}

func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := os.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}

type Data struct {
	Form *db.Recipe
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	_, err := db.InitDb()

	if err != nil {
		fmt.Println(err)
		return
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "./static/")

	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "pages/home/index.html", nil)
	})

	e.GET("/all", func(c echo.Context) error {
		log.Println("testing")
		return c.Render(200, "pages/all/index.html", nil)
	})
	e.POST("/recipes/add", routes.HandleRecipeAdd)
	e.GET("/recipes/all", routes.HandleGetAllRecipes)

	ip := fmt.Sprintf("%s:80", GetOutboundIP().String())
	log.Println(os.Getenv("ENV"))
	if os.Getenv("ENV") == "dev" {
		ip = "127.0.0.1:8080"
	}
	// Start server
	if err := e.Start(ip); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
