package data

import (
	"fmt"
	"sort"

	"github.com/IndependentIP/muse-scss-analyser/filesearch"
)

type Dependencies map[string][]filesearch.Match

type DependencyResult struct {
	File       string
	LineNumber uint32
	Symbol     string
	SourceFile string
}

type SearchResult interface {
	[]string | filesearch.Results | filesearch.FileStructure | Dependencies | []DependencyResult
}

func GetResults() filesearch.Results {
	return filesearch.MakeResultsFromResultsFile("generated/" + filesearch.LatestResultsFilename())
}

func FindAllImportsForFile(filename string, results *filesearch.Results) []string {
	fileData := (*results)[filename]

	return fileData.ImportedBy
}

func GetFileData(filename string, _ *filesearch.Results) filesearch.FileStructure {
	usages, err := filesearch.MakeFileStructureFromFilename(filename)

	if err != nil {
		fmt.Println("Error reading file: " + filename)
		fmt.Println(err)
		return filesearch.FileStructure{}
	}

	return usages
}

func FindUsedDependencies(file string, _ *filesearch.Results) []string {
	usages, err := filesearch.MakeFileStructureFromFilename(file)

	if err != nil {
		fmt.Println("Error reading file: " + file)
		fmt.Println(err)
		return []string{}
	}

	deps := make(map[string][]filesearch.Match)

	mergeDependencies(&deps, &usages.Functions)
	mergeDependencies(&deps, &usages.Variables)
	mergeDependencies(&deps, &usages.Mixins)

	results := make([]string, 0)

	for _, d := range deps {
		for _, m := range d {
			if !contains(results, m.File) {
				results = append(results, m.File)
			}
		}
	}
	return results
}

func FindUnusedDependencies(file string, results *filesearch.Results) []string {
	fileData := (*results)[file]
	imports := make([]string, 0)
	imports = append(imports, fileData.ImportedBy...)

	foundDependencies := FindUsedDependencies(file, results)

	return diffSlices(imports, foundDependencies)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func diffSlices(a, b []string) []string {
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		delete(m, item)
	}

	result := make([]string, 0)

	for item := range m {
		result = append(result, item)
	}

	return result
}

func MakeDependencyResultsFromDependencies(dependencies Dependencies) []DependencyResult {
	results := make([]DependencyResult, 0)

	for file, matches := range dependencies {
		for _, match := range matches {
			results = append(results, DependencyResult{
				File:       match.File,
				LineNumber: match.LineNumber,
				Symbol:     match.Match,
				SourceFile: file,
			})
		}
	}

	sortDependencyResults(&results)
	return results
}

func sortDependencyResults(deps *[]DependencyResult) {
	results := *deps

	sort.Slice(results, func(i, j int) bool {
		return results[i].LineNumber < results[j].LineNumber
	})
}

func FindAllDependenciesForFile(filename string, results *filesearch.Results) Dependencies {
	dependencies := make(Dependencies)
	fileData := (*results)[filename]
	imports := fileData.Matches

	for _, imp := range imports {
		if _, ok := dependencies[imp.File]; !ok {
			dependencies[imp.File] = make([]filesearch.Match, 0)
		}

		usages, err := filesearch.MakeFileStructureFromFilename(imp.File)

		if err != nil {
			fmt.Println("Error reading file: " + imp.File)
			fmt.Println(err)
			fmt.Println("Continuing...")
			continue
		}

		deps := make(map[string][]filesearch.Match)

		mergeDependencies(&deps, &usages.Functions)
		mergeDependencies(&deps, &usages.Variables)
		mergeDependencies(&deps, &usages.Mixins)

		fillDependenciesWithMatchesMap(filename, imp.File, deps, &dependencies)
	}

	return dependencies
}

func mergeDependencies(d1, d2 *map[string][]filesearch.Match) {
	for k, v := range *d2 {
		if (*d2)[k] == nil {
			continue
		}
		(*d1)[k] = v
	}
}

func fillDependenciesWithMatchesMap(
	filename string,
	importFile string,
	dependencyMatches map[string][]filesearch.Match,
	dependencies *Dependencies,
) {
	fmt.Println("Checking for dependencies in file: " + importFile)

	for _, d := range dependencyMatches {
		for _, m := range d {
			if m.File == filename {
				(*dependencies)[importFile] = append((*dependencies)[importFile], m)
			}
		}
	}
}
