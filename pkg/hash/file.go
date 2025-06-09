package hash

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const chunkSize = 65536 // 64k

func OsHashFromFilePath(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	hash, err := OsHashFromFile(f, fi.Size())
	sHash := fmt.Sprintf("%016x", hash)
	return sHash, nil
}

func OsHashFromFile(file *os.File, size int64) (string, error) {
	if size < chunkSize {
		return "", fmt.Errorf("file is too small")
	}

	buf := make([]byte, chunkSize*2)
	err := readChunk(file, 0, buf[:chunkSize])
	if err != nil {
		return "", err
	}
	err = readChunk(file, size-chunkSize, buf[chunkSize:])
	if err != nil {
		return "", err
	}

	var nums [(chunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return "", err
	}
	var hash uint64
	for _, num := range nums {
		hash += num
	}
	hash += uint64(size)

	return fmt.Sprintf("%016x", hash), nil
}

func readChunk(file *os.File, offset int64, buf []byte) error {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return err
	}
	if n != chunkSize {
		return fmt.Errorf("invalid read %v", n)
	}
	return nil
}
