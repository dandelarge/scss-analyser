package filesearch

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/tsx"
	scss "tree-sitter-scss"
)

type FileSearcher struct {
	typescriptParser *sitter.Parser
	scssParser       *sitter.Parser
}

func (fs *FileSearcher) ReadFileAndMakeTree(path string) (*sitter.Tree, *sitter.Language, []byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, err
	}

	isScss := strings.HasSuffix(path, ".scss")

	var parser *sitter.Parser
	var language *sitter.Language

	if isScss {
		parser = fs.scssParser
		language = scss.GetLanguage()
	} else {
		parser = fs.typescriptParser
		language = tsx.GetLanguage()
	}

	tree, err := parser.ParseCtx(context.Background(), nil, content)

	if err != nil {
		return nil, nil, nil, err
	}

	return tree, language, content, nil
}

func (fs *FileSearcher) GetFileMatches(file FileData) []Match {

	importStatementPattern := `(
		(import_statement) @import_statement
		(#match? @import_statement "\.scss")
	)`

	tree, language, content, err := fs.ReadFileAndMakeTree(file.path + "/" + file.file.Name())
	if err != nil {
		log.Fatal(err)
	}

	queryCursor := makeAndExecQuery(importStatementPattern, language, tree)

	matches := []Match{}

	for {
		match, ok := queryCursor.NextMatch()

		if !ok {
			break
		}

		match = queryCursor.FilterPredicates(match, content)

		for _, capture := range match.Captures {
			lineContent := capture.Node.Content(content)

			regexp := regexp.MustCompile(`['"](.*)['"]`)
			regexpMatch := regexp.FindStringSubmatch(lineContent)

			if len(regexpMatch) < 2 {
				log.Panic("No import file found in line" + ": " + lineContent)
			}

			if strings.HasPrefix(regexpMatch[1], "~") {
				continue
			}

			match := Match{
				File:       normalizePath(file.path + "/" + regexpMatch[1]),
				Line:       capture.Node.Content(content),
				LineNumber: capture.Node.StartPoint().Row + 1,
				Match:      normalizePath(file.path + "/" + regexpMatch[1]),
			}
			matches = append(matches, match)
		}
	}

	return matches
}

func (fs FileSearcher) FindPatternInFile(pattern string, path string) []Match {
	tree, language, content, err := fs.ReadFileAndMakeTree(path)
	if err != nil {
		log.Println(err)
		return []Match{}
	}

	queryCursor := makeAndExecQuery(pattern, language, tree)
	queryMatches := buildMatchesFromQueryCursor(queryCursor, content)

	matches := make([]Match, len(queryMatches))
	for _, match := range queryMatches {
		match.File = normalizePath(path)
		matches = append(matches, match)
	}

	return matches
}

func (fs FileSearcher) FindFunctionInvocations(sourceFile string) []Match {
	functionCallPattern := `
		(call_expression
			(function_name) @function_name)
	`

	return fs.FindPatternInFile(functionCallPattern, sourceFile)
}

func (fs FileSearcher) FindFunctionDeclarations(sourceFile string) []Match {
	functionStatementPattern := `
		(function_statement
			(name) @function)
	`

	return fs.FindPatternInFile(functionStatementPattern, sourceFile)
}

func (fs FileSearcher) FindIncludeStatements(sourceFile string) []Match {
	includeStatementPattern := `
		(include_statement
			(identifier) @include)
	`

	return fs.FindPatternInFile(includeStatementPattern, sourceFile)
}

func (fs FileSearcher) FindMixinDeclarations(sourceFile string) []Match {

	mixinStatementPattern := `
		(mixin_statement
			(name) @mixin)
	`

	return fs.FindPatternInFile(mixinStatementPattern, sourceFile)
}

func (fs FileSearcher) FindVarInvocations(sourceFile string) []Match {

	variableReferencePattern := `
		(variable_value) @variable
	`

	return fs.FindPatternInFile(variableReferencePattern, sourceFile)
}

func (fs FileSearcher) FindVarDeclarations(sourceFile string) []Match {

	variableStatementPattern := `
		(stylesheet
			(declaration
				(variable_name) @variable_name))
	`

	return fs.FindPatternInFile(variableStatementPattern, sourceFile)
}

func NewFileSearcher() *FileSearcher {

	typescriptParser := sitter.NewParser()
	scssParser := sitter.NewParser()

	typescriptParser.SetLanguage(tsx.GetLanguage())
	scssParser.SetLanguage(scss.GetLanguage())

	return &FileSearcher{
		typescriptParser: typescriptParser,
		scssParser:       scssParser,
	}
}

func makeAndExecQuery(queryPattern string, language *sitter.Language, tree *sitter.Tree) *sitter.QueryCursor {
	query, queryError := sitter.NewQuery([]byte(queryPattern), language)

	if queryError != nil {
		log.Fatal(queryError)
	}

	queryCursor := sitter.NewQueryCursor()

	queryCursor.Exec(query, tree.RootNode())
	return queryCursor
}

func buildMatchesFromQueryCursor(queryCursor *sitter.QueryCursor, content []byte) []Match {
	matches := []Match{}

	for {
		match, ok := queryCursor.NextMatch()

		if !ok {
			break
		}

		for _, capture := range match.Captures {
			lineContent := capture.Node.Content(content)

			lineMatch := Match{
				Line:       lineContent,
				LineNumber: capture.Node.StartPoint().Row + 1,
				Match:      lineContent,
			}

			matches = append(matches, lineMatch)
		}
	}

	return matches
}

func (fs *FileSearcher) FindInvocationsInResult(result Result, declarationsFn FindFunc, invocationsFn FindFunc, folderName string) map[string][]Match {
	fileName := result.File
	sourceFilePath := folderName + "/" + fileName
	declarations := declarationsFn(sourceFilePath)
	filesToSearch := result.ImportedBy
	matchMap := map[string][]Match{}

	for _, file := range filesToSearch {
		filePath := folderName + "/" + file
		fileExt := strings.Split(file, ".")[1]
		if fileExt != "scss" {
			continue
		}
		invocations := invocationsFn(filePath)
		for _, declaration := range declarations {
			currentMatches := matchMap[declaration.Match]
			if currentMatches == nil {
				currentMatches = []Match{}
			}
			matchMap[declaration.Match] = currentMatches
			for _, invocation := range invocations {
				if invocation.Match == declaration.Match && declaration.Match != "" {
					matchMap[declaration.Match] = append(matchMap[declaration.Match], invocation)
				}
			}
		}
	}

	return matchMap
}
