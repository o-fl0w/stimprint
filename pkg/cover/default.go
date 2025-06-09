package cover

import (
	"github.com/o-fl0w/stimprint/internal/ff"
)

var defaultParams = ff.Params{
	FrequencyLimit: 4000,
	FrequencyHints: []ff.FrequencyHint{
		{Hz: 2000, Color: "#00FF00"},
		{Hz: 800, Color: "#00FF00"},
		{Hz: 300, Color: "#FF7F00"},
		{Hz: 100, Color: "#FF0000"},
	},
	OutputImageWidth:       2400,
	OutputImageHeight:      480,
	WaveColorLeft:          "#1E90FF",
	WaveColorRight:         "#FF6347",
	WaveColorOverlap:       "#8A2BE2",
	WaveColorTriphase:      "#04FF04",
	WaveColorTriphaseAlpha: 0.7,
	WaveColorMono:          "#FFA500",
}

func DefaultParams() ff.Params {
	return defaultParams
}
