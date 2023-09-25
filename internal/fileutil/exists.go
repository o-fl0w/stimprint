package fileutil

import (
	"os"
)

func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
