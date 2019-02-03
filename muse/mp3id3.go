package muse

import (
	"github.com/bogem/id3v2"
	"strings"
)

type MP3ID3Handler struct {
}

func (p *MP3ID3Handler) MatchType() string {
	return "audio/mpeg"
}

func (p *MP3ID3Handler) ResultKey() string {
	return "id3"
}

func trimNullTerminator(s string) string {
	// Mp3tag has a bug where it adds NUL to our strings,
	// so strip them out as a compatibility hack.
	return strings.TrimSuffix(s, "\u0000")
}

func (p *MP3ID3Handler) Compute(filePath string) (interface{}, error) {
	mp3, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		return nil, err
	}
	defer mp3.Close()

	ret := map[string]string{}
	ret["title"] = trimNullTerminator(mp3.Title())
	ret["artist"] = trimNullTerminator(mp3.Artist())
	ret["album"] = trimNullTerminator(mp3.Album())
	ret["year"] = trimNullTerminator(mp3.Year())
	ret["genre"] = trimNullTerminator(mp3.Genre())
	return ret, nil
}
