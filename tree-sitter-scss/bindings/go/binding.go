package tree_sitter_scss

// #cgo CFLAGS: -std=c11 -fPIC
// #include "../../src/parser.c"
// #include "../../src/scanner.c"
import "C"

import (
	sitter "github.com/smacker/go-tree-sitter"
	"unsafe"
)

// Get the tree-sitter Language for this grammar.
func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_scss())
	return sitter.NewLanguage(ptr)
}
