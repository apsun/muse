package muse

import (
	"github.com/tcolgate/mp3"
	"io"
	"os"
	"time"
)

type MP3DurationHandler struct {
}

func (p *MP3DurationHandler) MatchType() string {
	return "audio/mpeg"
}

func (p *MP3DurationHandler) ResultKey() string {
	return "duration"
}

func (p *MP3DurationHandler) Compute(filePath string) (interface{}, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// We can't just guess the duration from the bitrate
	// because some files are variable-bitrate, and the
	// only way to distinguish CBR/VBR is to actually iterate
	// all the frames, so just always do it the slow way.
	totalDuration := time.Duration(0)
	d := mp3.NewDecoder(f)
	for {
		var frame mp3.Frame
		var skipped int
		err = d.Decode(&frame, &skipped)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		totalDuration += frame.Duration()
	}

	return int(totalDuration.Seconds()), nil
}
