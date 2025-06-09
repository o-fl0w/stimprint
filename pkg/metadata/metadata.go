package metadata

import (
	"context"
	"encoding/json"
	"github.com/o-fl0w/stimprint/internal/ff"
	"strconv"
	"time"
)

type probeResponse struct {
	Streams []struct {
		Channels int `json:"channels"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
		Tags     struct {
			Title  string `json:"title"`
			Artist string `json:"artist"`
		} `json:"tags"`
	} `json:"format"`
}

type Metadata struct {
	Channels int
	Duration time.Duration
	Title    string
	Artist   string
}

func GetMetadata(ctx context.Context, audioFilePath string) (Metadata, error) {
	args := []string{
		"-show_entries",
		"stream=channels:format=duration:format_tags=title,artist",
	}

	out, err := ff.ProbeOutput(ctx, audioFilePath, args)

	if err != nil {
		return Metadata{}, err
	}
	var r probeResponse
	err = json.Unmarshal(out, &r)
	if err != nil {
		return Metadata{}, err
	}
	fd, err := strconv.ParseFloat(r.Format.Duration, 32)
	if err != nil {
		return Metadata{}, err
	}

	md := Metadata{
		Channels: r.Streams[0].Channels,
		Duration: time.Duration(fd) * time.Second,
		Title:    r.Format.Tags.Title,
		Artist:   r.Format.Tags.Artist,
	}

	return md, nil
}
