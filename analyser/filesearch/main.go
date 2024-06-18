package filesearch

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func SearchFiles(folderName string) {
	fmt.Println("Searching for scss imports in folder: " + folderName)
	files, err := ReadRecursive(folderName)

	if err != nil {
		log.Fatal(err)
	}

	fs := NewFileSearcher()

	results := MakeResultsFromFileData(&files, fs)
	d3data := MakeD3Data(results)

	WriteResultsToFile(d3data, "generated/d3data.json")
	fmt.Println("Results for visualisation written to: generated/d3data.json")
	now := time.Now()
	formattedDate := fmt.Sprintf("__%d-%d-%d_%d-%d-%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	resultsFileName := fmt.Sprintf("generated/results%s.json", formattedDate)
	WriteResultsToFile(results, resultsFileName)

	fmt.Println("Finished finding imports!!")
	fmt.Println("Results written to: " + resultsFileName)
	fmt.Println("+--------------------------------+")
	fmt.Println("Finding invocations...")

	for _, result := range results {
		fileExtension := strings.Split(result.File, ".")[1]
		if fileExtension != "scss" {
			continue
		}

		fmt.Println("Processing file: " + result.File)

		varReferences := fs.FindInvocationsInResult(
			result,
			fs.FindVarDeclarations,
			fs.FindVarInvocations,
			folderName,
		)

		mixinReferences := fs.FindInvocationsInResult(
			result,
			fs.FindMixinDeclarations,
			fs.FindIncludeStatements,
			folderName,
		)

		funcReferences := fs.FindInvocationsInResult(
			result,
			fs.FindFunctionDeclarations,
			fs.FindFunctionInvocations,
			folderName,
		)

		fileStructure := MakeFileStructure()

		fileStructure.Variables = varReferences
		fileStructure.Mixins = mixinReferences
		fileStructure.Functions = funcReferences

		outputFileName := strings.ReplaceAll(result.File, "."+fileExtension, "_"+fileExtension+".json")
		outputFileName = strings.ReplaceAll(outputFileName, "/", "__")
		WriteResultsToFile(fileStructure, "generated/"+outputFileName)
	}
}
