package routes

import (
	"bufio"
	"fmt"
	"html/template"
	"meal-choices/db"
	"net/http"
	"strings"
)

func HandleCsvUpload(w http.ResponseWriter, r *http.Request, templates *template.Template) error {
	r.ParseMultipartForm(32 << 10)

	f, fHeader, err := r.FormFile("file")

	if f == nil || err != nil {
		return templates.ExecuteTemplate(w, "file-upload", "No file uploaded")
	}

	file, _ := fHeader.Open()
	defer file.Close()

	reader := bufio.NewReader(file)

	indexs := map[string]int{"page": -1, "book": -1, "name": -1}

	line, err := reader.ReadString('\n')
	headers := strings.Split(strings.ToLower(string(line[:len(line)-1])), ",")
	var invalidHeaders []string

	for i := range headers {
		header := strings.Trim(headers[i], " ")
		_, has := indexs[header]

		if !has {
			invalidHeaders = append(invalidHeaders, header)
		}

		indexs[header] = i
	}
	if len(invalidHeaders) > 0 {
		return templates.ExecuteTemplate(w, "file-upload", fmt.Sprintf("invalid headers: %v. \nPlease check the file and try again.", strings.Join(invalidHeaders, ", ")))
	}

	db := db.ConnectToDB()
	defer db.Close()

	line, err = reader.ReadString('\n')

	var count int64 = 0

	for string(line) != "" && err == nil {
		values := strings.Split(string(line[:len(line)-1]), ",")
		println(string(line))

		result, err := db.Exec(`INSERT INTO recipes (name, book, page) values (?,?,?)`, values[indexs["name"]], values[indexs["book"]], values[indexs["page"]])

		if err == nil {
			affected, err := result.RowsAffected()
			if err == nil {
				count += affected
			}
		}

		line, err = reader.ReadString('\n')
	}

	if err != nil {
		return err
	}

	return templates.ExecuteTemplate(w, "file-upload", fmt.Sprintf("Uploaded %v recipes successfully", count))

}
