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
		default:
			webapp.StartServer("4321")
		}
	}
}
