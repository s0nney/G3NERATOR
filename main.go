package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Article struct {
	Title      string
	Content    template.HTML
	Date       string
	Tags       []string
	MediaFiles []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bloggen <markdown-file.md>")
		os.Exit(1)
	}

	mdFile := os.Args[1]

	mdContent, err := ioutil.ReadFile(mdFile)
	if err != nil {
		log.Fatalf("Error reading markdown file: %v", err)
	}

	title, date, tags, content := parseMarkdown(mdContent)

	htmlContent := markdownToHTML(content)

	outputDir := slugify(title)
	if outputDir == "" {
		outputDir = "untitled-article"
	}

	mediaDir := filepath.Join(outputDir, "media")
	if err := os.MkdirAll(mediaDir, 0755); err != nil {
		log.Fatalf("Error creating media directory: %v", err)
	}

	imagePaths := findImagesInMarkdown(string(mdContent))
	var copiedImages []string
	for _, imgPath := range imagePaths {
		if strings.HasPrefix(imgPath, "http://") || strings.HasPrefix(imgPath, "https://") {
			continue
		}

		destPath := filepath.Join(mediaDir, filepath.Base(imgPath))
		if err := copyFile(imgPath, destPath); err == nil {
			copiedImages = append(copiedImages, destPath)
			htmlContent = strings.ReplaceAll(htmlContent,
				`src="`+imgPath+`"`,
				`src="media/`+filepath.Base(imgPath)+`"`)
		} else {
			log.Printf("Warning: Could not copy image %s: %v", imgPath, err)
		}
	}

	article := Article{
		Title:      title,
		Content:    template.HTML(htmlContent),
		Date:       date,
		Tags:       tags,
		MediaFiles: copiedImages,
	}

	copyStyleCSS(outputDir)

	generateHTML(outputDir, article)

	fmt.Printf("Successfully generated blog article in %s/\n", outputDir)
	fmt.Printf("Copied %d images to %s/media/\n", len(copiedImages), outputDir)
}

func slugify(s string) string {
	var result strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z':
			result.WriteRune(r)
		case r >= '0' && r <= '9':
			result.WriteRune(r)
		case r == ' ' || r == '-':
			result.WriteRune('-')
		}
	}
	return strings.Trim(result.String(), "-")
}

func findImagesInMarkdown(md string) []string {
	re := regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)|<img[^>]+src=["']([^"']+)["']`)
	matches := re.FindAllStringSubmatch(md, -1)

	var paths []string
	for _, match := range matches {
		if match[1] != "" {
			paths = append(paths, match[1])
		} else if match[2] != "" {
			paths = append(paths, match[2])
		}
	}
	return paths
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

func parseMarkdown(content []byte) (title, date string, tags []string, mdContent string) {
	lines := strings.Split(string(content), "\n")

	title = "Untitled"
	date = time.Now().Format("1/02/06")
	tags = []string{"untagged"}

	var contentLines []string
	inFrontMatter := false

	for _, line := range lines {
		if strings.HasPrefix(line, "---") {
			inFrontMatter = !inFrontMatter
			continue
		}

		if inFrontMatter {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "title":
					title = value
				case "date":
					date = value
				case "tags":
					tags = strings.Split(value, ",")
					for i, tag := range tags {
						tags[i] = strings.TrimSpace(tag)
					}
				}
			}
		} else {
			contentLines = append(contentLines, line)
		}
	}

	mdContent = strings.Join(contentLines, "\n")
	return
}

func markdownToHTML(md string) string {
	mdParser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := mdParser.Convert([]byte(md), &buf); err != nil {
		log.Fatalf("Error converting markdown to HTML: %v", err)
	}

	return buf.String()
}

func copyStyleCSS(outputDir string) {
	styleCSS := `@font-face {
  font-family: 'JetBrains Mono';
  src: url('fonts/JetBrainsMono-VariableFont_wght.ttf') format('truetype');
  font-weight: 100 800;
  font-style: normal;
  font-display: swap;
}

@font-face {
  font-family: 'JetBrains Mono Italic';
  src: url('fonts/JetBrainsMono-Italic-VariableFont_wght.ttf') format('truetype');
  font-weight: 100 800;
  font-style: italic;
  font-display: swap;
}

:root {
  --color-bg: black;
  --color-fg: #4d4d4d;
  --color-link: #626868;
  --color-link-visited: #FBFBFB;
  --color-link-hover: #ee6363;
}

html,
body {
  background: var(--color-bg);
  color: var(--color-fg);
  font-family: 'JetBrains Mono', 'JetBrains Mono Italic', monospace;
  height: 100%;
  width: 100%;
  margin: 0;
  padding: 0;
}

.container {
  display: grid;
  grid-template-columns: 1fr 460px 600px 1fr;
  grid-template-areas:
    ". left right .";
  column-gap: 80px;
  justify-items: center;
  align-items: center;
  min-height: 100%;
}

.left-container {
  grid-area: left;
  aspect-ratio: 1/1;
}

.right-container {
  grid-area: right;
  height: 50%;
  width: 100%;
}

.gif img {
  max-width: 100%;
  max-height: 100%;
}

.head {
  display: flex;
  flex-direction: column;
  font-size: 40px;
  padding-top: 60px;
}

.category {
  display: flex;
  flex-direction: column;
  width: 180px;
}

.main-box {
  display: flex;
  flex-direction: column;
  max-width: 750px;
  color: white;
  word-wrap: break-word;
}

.entries {
  border-radius: 2px;
  padding: 20px;
  border: 1px solid var(--color-fg);
  margin-bottom: 30px;
  background-color: rgba(77, 77, 77, 0.03);
  transition: all 0.2s ease;
}

.entries:hover {
  background-color: rgba(77, 77, 77, 0.08);
  border-color: var(--color-link);
}

.entries h1 {
  font-size: 24px;
  margin: 0 0 8px 0;
  color: var(--color-link-visited);
  letter-spacing: -0.5px;
}

.entries .date {
  color: var(--color-link);
  font-size: 14px;
  margin: 0 0 16px 0;
  display: block;
}

.entries p {
  margin: 0 0 16px 0;
  line-height: 1.6;
}

.entries .subtitle {
  color: var(--color-link);
  font-size: 14px;
  margin: 20px 0 0 0;
  border-top: 1px dashed var(--color-fg);
  padding-top: 12px;
}

.entries img {
  display: block;
  max-width: 100%;
  height: auto;
  margin: 20px auto;
  border: 1px solid var(--color-fg);
  padding: 4px;
  background-color: rgba(77, 77, 77, 0.03);
  transition: all 0.2s ease;
  filter: grayscale(20%);
}

.entries img:hover {
  border-color: var(--color-link);
  background-color: rgba(77, 77, 77, 0.08);
  filter: grayscale(0%);
  transform: scale(1.01);
}

.image-container {
  margin: 20px 0;
  text-align: center;
}

.image-caption {
  color: var(--color-link);
  font-size: 14px;
  margin-top: 8px;
  font-style: italic;
}

.tag {
  display: inline-block;
  background-color: rgba(77, 77, 77, 0.15);
  color: var(--color-link-visited);
  padding: 3px 8px;
  margin-right: 6px;
  margin-bottom: 4px;
  border-radius: 2px;
  font-size: 12px;
  text-transform: lowercase;
  letter-spacing: 0.5px;
  transition: all 0.2s ease;
}

.tag:hover {
  background-color: var(--color-link-hover);
  color: var(--color-bg);
  text-decoration: none;
}

.bookmarks {
  display: flex;
}

.links {
  display: flex;
  flex-direction: column;
  padding-top: 20px;
  padding-bottom: 20px;
}

.title {
  font-size: 20px;
}

li {
  font-size: 16px;
  list-style-type: none;
  padding: 5px;
}

a:link {
  text-decoration: none;
  color: var(--color-link);
}

a:visited {
  color: var(--color-link-visited);
}

a:hover {
  color: var(--color-link-hover);
}

.blinking {
  animation: opacity 1s ease-in-out infinite;
  opacity: 1;
}

@keyframes opacity {
  0% {
    opacity: 1;
  }

  50% {
    opacity: 0;
  }

  100% {
    opacity: 1;
  }
}

@media (max-width: 1024px) {
  .container {
    grid-template-columns: 1fr;
    grid-template-areas:
      "right";
    column-gap: 0;
  }

  .left-container {
    display: none;
  }

  .right-container {
    height: auto;
    width: 90%;
    padding: 20px;
    margin: 0 auto;
  }

  .head {
    padding-top: 40px;
    font-size: 32px;
  }

  .main-box {
    padding: 0 15px;
  }
}

@media (max-width: 480px) {
  .head {
    font-size: 24px;
    padding-top: 30px;
  }

  .category {
    width: 100%;
  }

  li {
    font-size: 14px;
  }

  .entries {
    padding: 15px;
  }

  .entries h1 {
    font-size: 20px;
  }

  .entries p {
    font-size: 15px;
  }

  .tag {
    font-size: 11px;
    padding: 2px 6px;
  }
}

footer {
  margin-top: 100px;
}`

	err := ioutil.WriteFile(filepath.Join(outputDir, "style.css"), []byte(styleCSS), 0644)
	if err != nil {
		log.Fatalf("Error writing style.css: %v", err)
	}
}

func generateHTML(outputDir string, article Article) {
	tmpl := `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" type="text/css" href="style.css">
    <link rel="icon" type="image/x-icon" href="favicon.ico">
</head>

<body>
    <div class="container">
        <div class="left-container">
            <div class="gif">
                <img src="media/banner.gif" />
            </div>
        </div>

        <div class="right-container">
            <div class="head">
                <p>entries<span class="blinking">_</span></p>
            </div>
            <div class="main-box">
                <article>
                    <div class="entries">
                        <h1>{{.Title}}</h1>
                        <span class="date">{{.Date}}</span>
                        {{.Content}}
                        <p class="subtitle">
                            {{range .Tags}}<span class="tag">{{.}}</span>{{end}}
                        </p>
                    </div>
                </article>
            </div>
        </div>
    </div>
</body>

</html>`

	t := template.Must(template.New("blog").Parse(tmpl))

	outputFile := filepath.Join(outputDir, "index.html")
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating HTML file: %v", err)
	}
	defer f.Close()

	data := struct {
		Title   string
		Date    string
		Content template.HTML
		Tags    []string
	}{
		Title:   article.Title,
		Date:    article.Date,
		Content: article.Content,
		Tags:    article.Tags,
	}

	if err := t.Execute(f, data); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}
