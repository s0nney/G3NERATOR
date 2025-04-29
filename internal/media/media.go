package media

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"g3nerator/internal/config"
)

type MediaProcessor struct {
	cfg *config.Config
}

func NewProcessor(cfg *config.Config) *MediaProcessor {
	return &MediaProcessor{cfg: cfg}
}

func (m *MediaProcessor) Process() (string, error) {
	content, err := os.ReadFile(m.cfg.MarkdownFile)
	if err != nil {
		return "", err
	}

	mdContent := string(content)

	// Create media directory
	if err := os.MkdirAll(m.cfg.MediaDir, 0755); err != nil {
		return "", fmt.Errorf("error creating media directory: %w", err)
	}

	// Process both markdown and HTML image references
	updatedContent := processImageReferences(mdContent, filepath.Dir(m.cfg.MarkdownFile), m.cfg.MediaDir)
	return updatedContent, nil
}

func processImageReferences(content, sourceDir, mediaDir string) string {
	// Regex for markdown images: ![alt](path)
	mdRe := regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
	// Regex for HTML images: <img src="path">
	htmlRe := regexp.MustCompile(`<img[^>]+src="([^">]+)"`)

	processMatch := func(match string, path string) string {
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			return match
		}

		absPath := path
		if !filepath.IsAbs(path) {
			absPath = filepath.Join(sourceDir, path)
		}

		fileName := filepath.Base(path)
		destPath := filepath.Join(mediaDir, fileName)

		if err := copyFile(absPath, destPath); err != nil {
			fmt.Printf("Warning: Could not copy media file %s: %v\n", path, err)
			return match
		}

		newPath := filepath.Join("media", fileName)
		return strings.Replace(match, path, newPath, 1)
	}

	content = mdRe.ReplaceAllStringFunc(content, func(match string) string {
		submatches := mdRe.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		return processMatch(match, submatches[1])
	})

	content = htmlRe.ReplaceAllStringFunc(content, func(match string) string {
		submatches := htmlRe.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		return processMatch(match, submatches[1])
	})

	return content
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
