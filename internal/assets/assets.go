package assets

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type AssetManager struct {
	SourceDir string
}

func NewAssetManager(sourceDir string) *AssetManager {
	return &AssetManager{
		SourceDir: sourceDir,
	}
}

func (a *AssetManager) CopyAssets(destDir string) error {
	assets := []string{
		"css/styles.css",
		"css/base.css",
		"css/layout.css",
		"css/components.css",
		"js/main.js",
		"js/animations.js",
	}

	for _, asset := range assets {
		srcPath := filepath.Join(a.SourceDir, asset)
		destPath := filepath.Join(destDir, asset)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create asset directory: %w", err)
		}

		if err := copyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to copy asset %s: %w", asset, err)
		}
	}

	return nil
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
