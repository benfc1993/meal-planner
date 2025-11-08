package pages

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func Handlers(f embed.FS) {

	templates, err := CreateTemplates(f, "public")

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Print(templates.DefinedTemplates())

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		var path string

		if len(r.URL.Path) == 1 {
			path = ""
		} else {
			path = r.URL.Path
		}

		dirPath := "public" + path
		_, err := f.ReadDir(dirPath)
		if err != nil {
			log.Print(err)
			w.WriteHeader(404)
			return
		}

		err = templates.ExecuteTemplate(w, "/", nil)

		if err != nil {
			log.Print(err)
			w.WriteHeader(404)
			return
		}
	})

	http.HandleFunc("GET /testing", func(w http.ResponseWriter, r *http.Request) {
		data := TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		}
		templates.ExecuteTemplate(w, "/testing", data)

	})
}

func CreateTemplates(f embed.FS, path string) (*template.Template, error) {

	root := template.New("")
	err := Walk(f, root, path, len(path))
	return root, err

}

func Walk(f embed.FS, root *template.Template, path string, prefix int) error {

	entries, err := f.ReadDir(path)
	if err != nil {
		return err
	}
	for i := range len(entries) {
		entry := entries[i]
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirPath := path + "/" + entry.Name()
			return Walk(f, root, dirPath, prefix)
		}
		if strings.HasSuffix(info.Name(), ".html") {
			filePath := path + "/" + info.Name()
			b, err := f.ReadFile(filePath)
			if err != nil {
				return err
			}
			if strings.HasPrefix(string(b), "<") {
				templateName := filePath[prefix : len(filePath)-len("page.html")]
				if strings.HasSuffix(templateName, "/") && len(templateName) > 1 {
					templateName = templateName[:len(templateName)-1]
				}
				t := root.New(templateName).Funcs(nil)
				t.Parse(string(b))
			} else {
				root.Parse(string(b))
			}
		}
	}
	return err
}
