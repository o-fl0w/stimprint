package ff

import (
	"fmt"
	"strconv"
	"strings"
)

type FrequencyHint struct {
	Hz    int
	Color string
}

func (fh *FrequencyHint) Set(s string) error {
	sfh := strings.Split(s, ":")
	hz, err := strconv.Atoi(sfh[0])
	if err != nil {
		return err
	}
	fh.Hz = hz
	fh.Color = sfh[1]
	return nil
}

func (fh *FrequencyHint) String() string {
	return fmt.Sprintf("%d:%s", fh.Hz, fh.Color)
}

func (fh *FrequencyHint) Y(frequencyLimit int, imageHeight int) int {
	pxPerHz := float32(imageHeight) / float32(frequencyLimit)
	return int(float32(imageHeight) - pxPerHz*float32(fh.Hz))
}

type FrequencyHints []FrequencyHint

func (v *FrequencyHints) Type() string {
	return "string"
}

func (v *FrequencyHints) Set(s string) error {
	ss := strings.Split(s, ";")
	fhs := make([]FrequencyHint, len(ss))
	for i := range fhs {
		err := fhs[i].Set(ss[i])
		if err != nil {
			return err
		}
	}
	*v = fhs
	return nil
}

func (v *FrequencyHints) String() string {
	ss := make([]string, len(*v))
	for i := range *v {
		ss[i] = (*v)[i].String()
	}
	return strings.Join(ss, ";")
}
