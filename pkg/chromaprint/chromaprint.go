//https://github.com/go-fingerprint/fingerprint/blob/master/fingerprint.go

package chromaprint

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"strconv"
	"strings"
)

const (
	bitsperint = 32
)

var (
	errLength = errors.New("fingerprints differs in length")
)

func FromBytes(bs []byte) ([]int32, error) {
	r := bytes.NewReader(bs)
	cp := make([]int32, len(bs)/4)
	err := binary.Read(r, binary.LittleEndian, &cp)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

// Compare returns a number that indicates how two fingerprints
// are similar to each other as a value from 0 to 1. Usually two
// fingerprints can be considered identical when the score is
// greater or equal than 0.95.
func Compare(fprint1, fprint2 []int32) (float64, error) {
	dist := 0
	if len(fprint1) != len(fprint2) {
		return 0, errLength
	}

	for i, sub := range fprint1 {
		dist += hamming(sub, fprint2[i])
	}

	score := 1 - float64(dist)/float64(len(fprint1)*bitsperint)
	return score, nil
}

// Distance returns slice of pairwisely XOR-ed fingerprints.
func Distance(fprint1, fprint2 []int32) ([]int32, error) {
	if len(fprint1) != len(fprint2) {
		return nil, errLength
	}

	dist := make([]int32, len(fprint1))

	for i, sub := range fprint1 {
		dist[i] = sub ^ fprint2[i]
	}
	return dist, nil
}

// ToImage returns black-and-white image.Image with graphical
// representation of fingerprint: each column represents
// 32-bit integer, where black and white pixels correspond
// to 1 and 0 respectively.
func ToImage(fprint []int32) (im image.Image) {
	return int32ToImage(fprint)
}

// ImageDistance returns black-and white image.Image with
// graphical representation of distance between fingerprints.
func ImageDistance(fprint1, fprint2 []int32) (im image.Image, err error) {
	if len(fprint1) != len(fprint1) {
		return nil, errLength
	}

	dist, err := Distance(fprint1, fprint2)
	if err != nil {
		return
	}
	im = int32ToImage(dist)
	return
}
func hamming(a, b int32) (dist int) {
	dist = strings.Count(strconv.FormatInt(int64(a^b), 2), "1")
	return
}

func int32ToImage(s []int32) image.Image {
	im := image.NewGray(image.Rect(0, 0, len(s), bitsperint))
	for i, sub := range s {
		for j := 0; j < bitsperint; j++ {
			im.Set(i, j, color.Gray{uint8(sub&1) * 0xFF})
			sub >>= 1
		}
	}
	return im
}
