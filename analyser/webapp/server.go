package webapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/IndependentIP/muse-scss-analyser/filesearch"
	"github.com/IndependentIP/muse-scss-analyser/webapp/data"
	templates "github.com/IndependentIP/muse-scss-analyser/webapp/templates"
	"github.com/gorilla/mux"
)

const testFile = "styles/responsive.scss"

func StartServer(port string) {

	r := mux.NewRouter()

	r.HandleFunc("/api/results", GetResults).Methods("GET")

	r.HandleFunc("/api/dependencies", GetDependencies).Methods("GET")

	GetUnused := GetHandlerForFunction(data.FindUnusedDependencies)
	r.HandleFunc("/api/unused", GetUnused).Methods("GET")

	GetUsed := GetHandlerForFunction(data.FindUsedDependencies)
	r.HandleFunc("/api/used", GetUsed).Methods("GET")

	GetImports := GetHandlerForFunction(data.FindAllImportsForFile)
	r.HandleFunc("/api/imports", GetImports).Methods("GET")

	GetFileData := GetHandlerForFunction(data.GetFileData)
	r.HandleFunc("/api/filedata", GetFileData).Methods("GET")

	fs := http.FileServer(http.Dir("./webapp/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Index().Render(r.Context(), w)
	})

	fmt.Println("Server started at http://localhost:" + port)

	error := http.ListenAndServe(":"+port, r)

	if error != nil {
		log.Fatal("Server failed to start: ", error)
	}
}

func GetResults(w http.ResponseWriter, r *http.Request) {
	data := data.GetResults()

	resultsJson, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resultsJson)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func GetDependencies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	file := query.Get("f")
	results := data.GetResults()

	if file == "" {
		filenames := make([]string, 0, len(results))

		for filename := range results {
			filenames = append(filenames, filename)
		}

		filenamesJson, err := json.Marshal(filenames)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(filenamesJson)
		return
	}

	decodedFilename, err := url.QueryUnescape(file)

	fmt.Println("File: ", decodedFilename)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dependencies := data.FindAllDependenciesForFile(decodedFilename, &results)
	dependenciesJson, err := json.Marshal(dependencies)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(dependenciesJson)
}

func GetHandlerForFunction[T data.SearchResult](f func(file string, results *filesearch.Results) T) func(path http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		file := query.Get("f")
		decodedFilename, err := url.QueryUnescape(file)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		results := data.GetResults()
		dependencies := f(decodedFilename, &results)

		dependenciesJson, err := json.Marshal(dependencies)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(dependenciesJson)
	}
}
