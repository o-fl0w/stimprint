package fileutil

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

const prefix = "stimhub"

func MkTmpFilePath(id uint64) string {
	s := filepath.Join(os.TempDir(), prefix+strconv.FormatUint(id, 10))
	return s
}

func MkRandomTmpFilePath() string {
	return MkTmpFilePath(rand.Uint64())
}
