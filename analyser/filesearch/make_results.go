package filesearch

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func MakeFileStructureFromFilename(filename string) (FileStructure, error) {

	matchFilename := strings.ReplaceAll(filename, "/", "__")
	matchFilename = strings.ReplaceAll(matchFilename, ".", "_")
	matchFilename = "generated/" + matchFilename + ".json"

	fmt.Println("Reading file: " + matchFilename)

	file, err := os.ReadFile(matchFilename)
	if err != nil {
		return FileStructure{}, err
	}

	var fileStructure FileStructure
	err = json.Unmarshal(file, &fileStructure)

	if err != nil {
		return FileStructure{}, err
	}

	return fileStructure, nil
}

func MakeResultsFromResultsFile(resultsFile string) Results {
	file, err := os.ReadFile(resultsFile)
	if err != nil {
		log.Fatal(err)
	}

	var results Results
	err = json.Unmarshal(file, &results)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

func MakeResultsFromFileData(fileData *[]FileData, fileSearcher *FileSearcher) Results {
	results := Results{}

	for _, file := range *fileData {
		// read file contents
		matches := fileSearcher.GetFileMatches(file)
		if len(matches) == 0 {
			continue
		}
		fmt.Println("Matches found in file: " + file.file.Name())

		fileName := normalizePath(file.path + "/" + file.file.Name())

		var fileResult Result

		if _, ok := results[fileName]; ok {
			fileResult = results[fileName]

		} else {
			fileResult = Result{
				File:       fileName,
				ImportedBy: []string{},
			}
		}

		fileResult.Matches = matches
		results[fileName] = fileResult

		for _, match := range matches {
			importedBy, ok := results[match.File]

			if !ok {
				results[match.File] = Result{
					File:       match.File,
					Matches:    []Match{},
					ImportedBy: []string{fileName},
				}
				continue
			}

			importedBy.ImportedBy = append(importedBy.ImportedBy, fileName)
			results[match.File] = importedBy
		}

	}

	return results
}

func MakeD3Data(results Results) D3Data {
	nodes := []Node{}
	links := []Link{}

	for _, result := range results {
		nodes = append(nodes, Node{
			Id:   result.File,
			Name: result.File,
		})

		for _, importedBy := range result.ImportedBy {
			links = append(links, Link{
				Source: importedBy,
				Target: result.File,
			})
		}
	}

	return D3Data{
		Nodes: nodes,
		Links: links,
	}
}

func MakeFileStructure() FileStructure {
	return FileStructure{
		Variables: map[string][]Match{},
		Functions: map[string][]Match{},
		Mixins:    map[string][]Match{},
	}
}

func WriteResultsToFile(results any, fileName string) {
	resultsJson, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	targetFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer targetFile.Close()

	targetFile.WriteString(string(resultsJson))

	log.Println("Results written to: " + fileName)
}

func sortAlphabetically(files []string) []string {
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			if files[i] > files[j] {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
	return files
}

func ListResultsFiles() []string {
	resultsFiles := []string{}

	files, err := os.ReadDir("generated")
	if err != nil {
		fmt.Println("Error reading the 'generated' folder")
		panic(err)
	}

	for _, file := range files {
		if strings.Index(file.Name(), "results") > -1 {
			resultsFiles = append(resultsFiles, file.Name())
		}
	}
	return sortAlphabetically(resultsFiles)
}

func LatestResultsFilename() string {
	resultsFiles := ListResultsFiles()
	lastResultFile := resultsFiles[len(resultsFiles)-1:][0]
	return lastResultFile
}
