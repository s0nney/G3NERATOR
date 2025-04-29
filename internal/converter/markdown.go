package converter

import (
	"bytes"

	"g3nerator/internal/config"

	"github.com/yuin/goldmark"
)

type MarkdownConverter struct {
	cfg *config.Config
}

func NewMarkdownConverter(cfg *config.Config) *MarkdownConverter {
	return &MarkdownConverter{cfg: cfg}
}

func (m *MarkdownConverter) ConvertToHTML(mdContent string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(mdContent), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
