package param

type Cover struct {
	FrequencyLimit         int
	FrequencyHints         []FrequencyHint
	OutputImageWidth       int
	OutputImageHeight      int
	WaveColorLeft          string
	WaveColorRight         string
	WaveColorOverlap       string
	WaveColorTriphase      string
	WaveColorTriphaseAlpha float64
	WaveColorMono          string
}

type FrequencyHint struct {
	Hz    int
	Color string
}

var defaultParams = Cover{
	FrequencyLimit: 4000,
	FrequencyHints: []FrequencyHint{
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

func DefaultCoverParams() Cover {
	return defaultParams
}
