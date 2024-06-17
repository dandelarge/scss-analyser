package webapp

import (
	"fmt"
	templates "github.com/IndependentIP/muse-scss-analyser/webapp/templates"
	"github.com/a-h/templ"
	"net/http"
)

func StartServer(port string) {
	component := templates.Index()
	fs := http.FileServer(http.Dir("./webapp/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", templ.Handler(component))

	fmt.Println("Server started at http://localhost:" + port)

	error := http.ListenAndServe(":"+port, nil)

	if error != nil {
		fmt.Println("Error starting server: ", error)
	}
}
