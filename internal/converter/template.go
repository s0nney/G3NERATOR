package converter

import (
	"html/template"
	"os"
	"path/filepath"

	"g3nerator/internal/config"
)

type BlogPost struct {
	Title    string
	Content  template.HTML
	Date     string
	Author   string
	Category string
	ReadTime string
	Tags     []string
	MediaDir string
}

func firstChar(s string) string {
	if len(s) > 0 {
		return string(s[0])
	}
	return ""
}

func GenerateHTMLFile(cfg *config.Config, postData map[string]interface{}, htmlContent string) error {
	post := BlogPost{
		Title:    postData["title"].(string),
		Content:  template.HTML(htmlContent),
		Date:     postData["date"].(string),
		Author:   postData["author"].(string),
		Category: postData["category"].(string),
		ReadTime: postData["readtime"].(string),
		Tags:     postData["tags"].([]string),
		MediaDir: cfg.OutputDir,
	}

	tmpl := template.New("base.tmpl").Funcs(template.FuncMap{
		"firstChar": firstChar,
	})

	tmpl, err := tmpl.ParseFiles(
		"templates/base.html",
		"templates/article.html",
		"templates/partials/header.html",
		"templates/partials/footer.html",
		"templates/partials/sidebar.html",
	)
	if err != nil {
		return err
	}

	htmlPath := filepath.Join(cfg.OutputDir, cfg.HTMLName+".html")
	file, err := os.Create(htmlPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.ExecuteTemplate(file, "base.html", post)
}
