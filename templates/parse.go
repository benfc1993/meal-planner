package templates

import (
	"embed"
	"html/template"
	"log"
	"strings"
)

func CreateTemplates(f embed.FS, path string) (*template.Template, error) {

	root := template.New("")
	err := Walk(f, root, path, len(path))
	log.Print(root.DefinedTemplates())
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
			Walk(f, root, dirPath, prefix)
		}
		if strings.HasSuffix(info.Name(), ".html") || strings.HasSuffix(info.Name(), ".tmpl") {
			filePath := path + "/" + info.Name()
			b, err := f.ReadFile(filePath)
			if err != nil {
				return err
			}
			if info.Name() == "page.html" {
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
