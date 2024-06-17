package filesearch

import (
	"strings"
)

func normalizePath(path string) string {
	sections := strings.Split(path, "/")
	var normalizedSections []string
	var warpedSections []string
	var currentFolderSections []string
	var relativePathSections []string
	warpIdxStart := -1
	warpIdxEnd := -1

	for idx, section := range sections {
		if section == "." {
			warpIdxStart = idx
			warpIdxEnd = idx
			continue
		}
		if section == ".." {
			if warpIdxStart == -1 {
				warpIdxStart = idx
			}
			warpIdxEnd = idx
			continue
		}
	}

	if warpIdxStart == -1 {
		return strings.Join(sections[6:], "/")
	}

	warpedSections = sections[warpIdxStart : warpIdxEnd+1]
	var warpSize int

	if warpedSections[0] == "." {
		warpSize = len(warpedSections) - 1
	} else {
		warpSize = len(warpedSections)
	}
	currentFolderSections = sections[6 : warpIdxStart-warpSize]
	relativePathSections = sections[warpIdxEnd+1:]
	normalizedSections = append(currentFolderSections, relativePathSections...)

	return strings.Join(normalizedSections, "/")
}
