package api

import (
	"fmt"
	"net/http"
)

func Handlers() {
	h1 := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling /testing")
		fmt.Fprintf(w, "Hello")
	}

	http.HandleFunc("GET /api/testing", h1)
}
