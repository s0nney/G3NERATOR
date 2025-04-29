package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"g3nerator/internal/assets"
	"g3nerator/internal/config"
	"g3nerator/internal/converter"
	"g3nerator/internal/media"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <markdown-file>")
		return
	}

	markdownFile := os.Args[1]
	cfg := config.LoadConfig(markdownFile)

	mediaProcessor := media.NewProcessor(cfg)
	mdConverter := converter.NewMarkdownConverter(cfg)
	assetManager := assets.NewAssetManager(cfg.AssetsDir)

	if err := os.MkdirAll(filepath.Join(cfg.OutputDir, "css"), 0755); err != nil {
		log.Fatalf("Error creating css directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(cfg.OutputDir, "js"), 0755); err != nil {
		log.Fatalf("Error creating js directory: %v", err)
	}

	updatedContent, err := mediaProcessor.Process()
	if err != nil {
		log.Fatalf("Error processing media: %v", err)
	}

	if err := assetManager.CopyAssets(cfg.OutputDir); err != nil {
		log.Printf("Warning: could not copy all assets: %v", err)
	}

	htmlContent, err := mdConverter.ConvertToHTML(updatedContent)
	if err != nil {
		log.Fatalf("Error converting markdown: %v", err)
	}

	postData := map[string]interface{}{
		"title":    extractTitle(updatedContent),
		"date":     time.Now().Format("Jan 02, 2006"),
		"author":   "OPERATOR",
		"category": "SECURITY RESEARCH",
		"readtime": calculateReadTime(updatedContent),
		"tags":     []string{"EXPLOIT", "SECURITY"},
	}

	if err := converter.GenerateHTMLFile(cfg, postData, htmlContent); err != nil {
		log.Fatalf("Error generating HTML: %v", err)
	}

	fmt.Printf("Successfully created blog post in directory %s\n", cfg.OutputDir)
}

func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return "Untitled Post"
}

func calculateReadTime(content string) string {
	wordCount := len(strings.Fields(content))
	minutes := wordCount / 200
	if minutes < 1 {
		return "1"
	}
	return fmt.Sprintf("%d", minutes)
}
