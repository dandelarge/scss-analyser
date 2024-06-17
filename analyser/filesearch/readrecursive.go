package filesearch

import (
	"os"
	"path/filepath"
)

func filterFiles(files []FileData) []FileData {
	var filteredFiles []FileData

	fileExtensions := map[string]bool{
		".js":   true,
		".jsx":  true,
		".ts":   true,
		".tsx":  true,
		".scss": true,
	}

	for _, file := range files {
		extension := filepath.Ext(file.file.Name())

		if !file.file.IsDir() && fileExtensions[extension] {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

func ReadRecursive(folderName string) ([]FileData, error) {
	files, err := os.ReadDir(folderName)
	if err != nil {
		return nil, err
	}

	var fileData []FileData

	for _, file := range files {
		if file.IsDir() {
			subfolderFiles, err := ReadRecursive(folderName + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			fileData = append(fileData, subfolderFiles...)
		} else {
			fileData = append(fileData, FileData{
				file: file,
				path: folderName,
			})
		}
	}

	return filterFiles(fileData), nil
}
