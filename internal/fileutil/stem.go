package fileutil

import "path/filepath"

func Stem(filePath string) string {
	fileName := filepath.Base(filePath)
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
