package muse

import (
	"os"
	"path/filepath"
)

type PreviewHandler interface {
	MatchType() string
	ResultKey() string
	Compute(filePath string) (interface{}, error)
}

type PreviewProvider struct {
	handlers map[string]map[string][]PreviewHandler
}

func (s *PreviewProvider) getPreviewHandlers(mimeType string) []PreviewHandler {
	// Match */*, major/*, and major/minor
	major, minor := splitMIMEType(mimeType)
	ret := []PreviewHandler{}
	ret = append(ret, s.handlers["*"]["*"]...)
	ret = append(ret, s.handlers[major]["*"]...)
	ret = append(ret, s.handlers[major][minor]...)
	return ret
}

func (s *PreviewProvider) previewFile(filePath string, fileInfo os.FileInfo) map[string]interface{} {
	ret := map[string]interface{}{}
	ret["name"] = fileInfo.Name()
	if fileInfo.IsDir() {
		ret["type"] = "directory"
	} else {
		mimeType := guessMIMEType(filePath)
		ret["type"] = mimeType
		for _, handler := range s.getPreviewHandlers(mimeType) {
			prop, err := handler.Compute(filePath)
			if err == nil {
				ret[handler.ResultKey()] = prop
			}
		}
	}
	return ret
}

func (s *PreviewProvider) previewDirectory(dirPath string, dir *os.File) ([]map[string]interface{}, error) {
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	ret := make([]map[string]interface{}, len(fileInfos))
	for i, fileInfo := range fileInfos {
		filePath := filepath.Join(dirPath, fileInfo.Name())
		preview := s.previewFile(filePath, fileInfo)
		ret[i] = preview
	}

	return ret, nil
}

func (s *PreviewProvider) Register(handler PreviewHandler) {
	matchType := handler.MatchType()
	major, minor := splitMIMEType(matchType)
	minorHandlers := s.handlers[major]
	if minorHandlers == nil {
		minorHandlers = map[string][]PreviewHandler{}
		s.handlers[major] = minorHandlers
	}
	s.handlers[major][minor] = append(s.handlers[major][minor], handler)
}

func NewPreviewProvider() *PreviewProvider {
	return &PreviewProvider{
		handlers: map[string]map[string][]PreviewHandler{},
	}
}
