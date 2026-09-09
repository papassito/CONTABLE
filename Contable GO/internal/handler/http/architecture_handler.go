package http

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GoFileDTO represents the structure of a source code file for the frontend.
type GoFileDTO struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Layer       string `json:"layer"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type ArchitectureHandler struct {
	rootPath string
}

// NewArchitectureHandler creates a new instance of the inspection handler.
func NewArchitectureHandler(rootPath string) *ArchitectureHandler {
	return &ArchitectureHandler{rootPath: rootPath}
}

// GetFileTree inspects the local disk and serves the current structure as JSON.
func (h *ArchitectureHandler) GetFileTree(w http.ResponseWriter, r *http.Request) {
	var files []GoFileDTO

	err := filepath.WalkDir(h.rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		// We only filter for Go code files and JSON schemas
		if !strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".json") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(h.rootPath, path)
		layer := determineLayer(relPath)

		files = append(files, GoFileDTO{
			Path:        filepath.ToSlash(relPath),
			Name:        d.Name(),
			Layer:       layer,
			Description: "FCOS v2.2 architecture source file",
			Content:     string(content),
		})
		return nil
	})

	if err != nil {
		http.Error(w, "Error inspecting file system", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// determineLayer classifies each file within the Clean Architecture layers.
func determineLayer(relPath string) string {
	cleanPath := filepath.ToSlash(relPath)
	switch {
	case strings.HasPrefix(cleanPath, "cmd/"):
		return "Entrypoint"
	case strings.HasPrefix(cleanPath, "internal/domain") || strings.Contains(cleanPath, "/domain/"):
		return "Domain"
	case strings.HasPrefix(cleanPath, "internal/service") || strings.Contains(cleanPath, "/application/"):
		return "Application"
	case strings.HasPrefix(cleanPath, "internal/repository") || strings.Contains(cleanPath, "/infrastructure/"):
		return "Infrastructure"
	case strings.HasPrefix(cleanPath, "internal/handler"):
		return "Interface"
	case strings.HasPrefix(cleanPath, "pkg/"):
		return "Shared Core"
	default:
		return "Configuration"
	}
}
