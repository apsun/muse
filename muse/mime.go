package muse

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func trimMIMEType(mime string) string {
	end := strings.Index(mime, ";")
	if end >= 0 {
		mime = mime[:end]
	}
	return mime
}

func splitMIMEType(mime string) (string, string) {
	parts := strings.Split(mime, "/")
	if len(parts) != 2 {
		panic("Cannot parse MIME type")
	}
	return parts[0], parts[1]
}

func guessMIMEType(filePath string) string {
	// If we can figure out the type by looking at the file extension,
	// use that information.
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType != "" {
		return trimMIMEType(mimeType)
	}

	// File has no extension, so guess the MIME type by looking at the
	// file header.
	f, err := os.Open(filePath)
	if err != nil {
		return "application/octet-stream"
	}
	defer f.Close()

	// Read as much as possible in one call. We don't want to hold up
	// the rest of the code by waiting for data, so just pray that
	// enough data is returned in the first call.
	head := make([]byte, 512)
	n, err := f.Read(head)
	if err != nil {
		return "application/octet-stream"
	}
	head = head[:n]

	mimeType = http.DetectContentType(head)
	return trimMIMEType(mimeType)
}
