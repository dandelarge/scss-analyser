package main

import (
	"fmt"
	"os"

	"github.com/IndependentIP/muse-scss-analyser/filesearch"
	"github.com/IndependentIP/muse-scss-analyser/webapp"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		webapp.StartServer("4321")
	} else {
		command := args[0]
		switch command {
		case "inspect":
		case "i":
		case "find":
		case "search":
			if len(args) < 2 {
				fmt.Println("Please provide a folder to search")
				os.Exit(1)
			}
			folderName := args[1]
			filesearch.SearchFiles(folderName)
		case "serve":
			var port string
			if len(args) < 2 {
				port = "4321"
			} else {
				port = args[1]
			}
			webapp.StartServer(port)
		case "gr":
		case "graph-results":
			var fileName string
			if len(args) < 2 {
				fileName = filesearch.LatestResultsFilename()
			} else {
				fileName = args[1]
			}

			results := filesearch.MakeResultsFromResultsFile("generated/" + fileName)
			d3data := filesearch.MakeD3Data(results)
			filesearch.WriteResultsToFile(d3data, "generated/d3data.json")
		default:
			webapp.StartServer("4321")
		}
	}
}
