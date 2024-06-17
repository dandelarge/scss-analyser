package webapp

import (
	"fmt"
	templates "github.com/IndependentIP/muse-scss-analyser/webapp/templates"
	"github.com/a-h/templ"
	"net/http"
)

func StartServer() {
	component := templates.Index()
	fs := http.FileServer(http.Dir("webapp/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :4321")

	http.ListenAndServe(":4321", nil)
}
