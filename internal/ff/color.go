package ff

import "strconv"

type rgb struct {
	r uint8
	g uint8
	b uint8
}

func hex2rgb(hex string) rgb {
	if hex[0] == '#' {
		hex = hex[1:]
	}
	var c rgb
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return rgb{}
	}

	c = rgb{
		r: uint8(values >> 16),
		g: uint8((values >> 8) & 0xFF),
		b: uint8(values & 0xFF),
	}

	return c
}
