package muse

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Server struct {
	root     string
	provider *PreviewProvider
}

func (s *Server) getFilesystemPath(urlPath string) string {
	root := s.root
	if root == "" {
		root = "."
	}
	return filepath.Join(root, filepath.FromSlash(path.Clean(urlPath)))
}

func (s *Server) serveDirectory(w http.ResponseWriter, dirPath string, dir *os.File) {
	// Generate response
	result, err := s.provider.previewDirectory(dirPath, dir)
	if err != nil {
		log.Printf("Failed to generate dir preview: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write response body. If this fails, there's nothing we can do
	// since we already wrote the header. Just ignore this and let the
	// client deal with it.
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("Failed to encode JSON: %s\n", err.Error())
		return
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := s.getFilesystemPath(r.URL.Path)

	f, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open path: %s\n", err.Error())
		if os.IsNotExist(err) {
			http.Error(w, "Path does not exist", http.StatusNotFound)
		} else if os.IsPermission(err) {
			http.Error(w, "Cannot access path", http.StatusForbidden)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		log.Printf("Failed to stat path: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		s.serveDirectory(w, path, f)
	} else {
		http.ServeContent(w, r, path, info.ModTime(), f)
	}
}

func (s *Server) RegisterPreviewHandler(handler PreviewHandler) {
	s.provider.Register(handler)
}

func NewServer(root string) *Server {
	return &Server{
		root:     root,
		provider: NewPreviewProvider(),
	}
}
