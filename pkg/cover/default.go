package cover

import (
	"github.com/o-fl0w/stimprint/internal/ff"
)

var defaultParams = Params{
	Ffprobe:                "ffprobe",
	Ffmpeg:                 "ffmpeg",
	OverwriteExistingFiles: false,
	FrequencyLimit:         4000,
	FrequencyHints: []ff.FrequencyHint{
		{Hz: 2000, Color: "#00FF00"},
		{Hz: 800, Color: "#00FF00"},
		{Hz: 300, Color: "#FF7F00"},
		{Hz: 100, Color: "#FF0000"},
	},
	OutputImageWidth:       1200,
	OutputImageHeight:      600,
	WaveColorLeft:          "#1E90FF",
	WaveColorRight:         "#FF6347",
	WaveColorOverlap:       "#EFCFEF",
	WaveColorTriphase:      "#04FF04",
	WaveColorTriphaseAlpha: 0.6,
	WaveColorMono:          "#FFA500",
}
