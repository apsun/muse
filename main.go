package main

import (
	"./muse"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	rootDir := flag.String("root", ".", "Root server directory")
	port := flag.Uint("port", 8080, "Port to listen on")
	flag.Parse()

	s := muse.NewServer(*rootDir)
	s.RegisterPreviewHandler(&muse.MP3ID3Handler{})

	// TODO: Ridiculously slow, need caching support before enabling
	// s.RegisterPreviewHandler(&metaserver.MP3DurationHandler{})

	addr := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(addr, s)
	if err != nil {
		panic(err.Error())
	}
}
