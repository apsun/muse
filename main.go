package main

import (
	"./muse"
	"net/http"
)

func main() {
	s := muse.NewServer("/home/andrew/Music")
	s.RegisterPreviewHandler(&muse.MP3ID3Handler{})

	// TODO: Ridiculously slow, need caching support before enabling
	// s.RegisterPreviewHandler(&metaserver.MP3DurationHandler{})

	http.ListenAndServe(":8080", s)
}
