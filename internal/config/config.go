package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	MarkdownFile string
	OutputDir    string
	MediaDir     string
	HTMLName     string
	AssetsDir    string
}

func LoadConfig(markdownFile string) *Config {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter HTML filename (without extension): ")
	htmlName, _ := reader.ReadString('\n')
	htmlName = strings.TrimSpace(htmlName)
	if htmlName == "" {
		htmlName = "index"
	}

	dirNumber := findNextAvailableDir()
	dirName := fmt.Sprintf("%03d", dirNumber)
	mediaDir := filepath.Join(dirName, "media")

	return &Config{
		MarkdownFile: markdownFile,
		OutputDir:    dirName,
		MediaDir:     mediaDir,
		HTMLName:     htmlName,
		AssetsDir:    "assets",
	}
}

func findNextAvailableDir() int {
	dirNumber := 1
	for {
		dirPath := fmt.Sprintf("%03d", dirNumber)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			break
		}
		dirNumber++
	}
	return dirNumber
}
