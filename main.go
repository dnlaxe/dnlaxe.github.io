package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"go.yaml.in/yaml/v3"
)

type PostData struct {
	Title           string
	Date            string
	Tags            []string
	Link            string
	Content         template.HTML
	GitHubUrl       string
	AdventOfCodeUrl string
	LiveDemo        string
}

type postFrontmatter struct {
	Title           string
	Date            string
	Type            string
	Slug            string
	Tags            []string
	GitHubUrl       string
	AdventOfCodeUrl string
	LiveDemo        string
}

type BlogPostData struct {
	Tags  []string
	Posts []PostData
}

func main() {

	serveFlag := flag.Bool("serve", false, "Start a local web server after building.")
	flag.Parse()

	os.RemoveAll("public")

	layoutTmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/navbar.html", "templates/themeToggle.html"))
	homeTmpl := template.Must(template.ParseFiles("templates/home.html", "templates/navbar.html", "templates/themeToggle.html"))
	blogTmpl := template.Must(template.ParseFiles("templates/blog.html", "templates/navbar.html", "templates/themeToggle.html"))
	portfolioTmpl := template.Must(template.ParseFiles("templates/portfolio.html", "templates/navbar.html", "templates/themeToggle.html"))
	notFoundTmpl := template.Must(template.ParseFiles("templates/404.html", "templates/navbar.html", "templates/themeToggle.html"))

	var posts []PostData
	var projects []PostData

	files, err := os.ReadDir("content")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}

		post, postType := parseMarkdown(filepath.Join("content", file.Name()), file.Name())

		switch postType {
		case "post":
			post.Link = "/blog/" + post.Link
			posts = append(posts, post)
		case "project":
			post.Link = "/portfolio/" + post.Link
			projects = append(projects, post)
		}

		folderName := strings.TrimPrefix(post.Link, "/")
		folderName = strings.TrimSuffix(folderName, "/")
		generateHTML(layoutTmpl, post, folderName, "index.html")

	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date > posts[j].Date
	})

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Date > projects[j].Date
	})

	fmt.Println("Building homepage...")
	generateHTML(homeTmpl, nil, "", "index.html")

	fmt.Println("Building blog...")

	uniqueTagsMap := make(map[string]bool)
	for i := range posts {
		for j, tag := range posts[i].Tags {
			cleanedTag := strings.ToUpper(strings.TrimSpace(tag))
			posts[i].Tags[j] = cleanedTag
			uniqueTagsMap[cleanedTag] = true
		}
	}

	var allTags []string
	for tag := range uniqueTagsMap {
		allTags = append(allTags, tag)
	}

	blogData := BlogPostData{
		Tags:  allTags,
		Posts: posts,
	}

	generateHTML(blogTmpl, blogData, "blog", "index.html")

	fmt.Println("Building portfolio...")
	generateHTML(portfolioTmpl, projects, "portfolio", "index.html")

	fmt.Println("Building 404 page...")
	generateHTML(notFoundTmpl, nil, "", "404.html")

	fmt.Println("Copying static assets...")
	copyStaticAssets()

	fmt.Println("Site fully generated.")

	if *serveFlag {
		fmt.Println("Starting local server! Open your browser to http://localhost:8080")
		fmt.Println("(Press CTRL+C in this terminal to stop the server)")
		log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("public"))))
	} else {
		fmt.Println("Build finished. Exiting process.")
	}
}

func parseMarkdown(filePath string, filename string) (PostData, string) {
	fmt.Printf("Processing: %s...\n", filename)

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal()
	}

	parts := strings.SplitN(string(fileBytes), "---", 3)
	if len(parts) < 3 {
		log.Fatalf("Error: File %s is missing formatter", filename)
	}

	var meta postFrontmatter
	err = yaml.Unmarshal([]byte(parts[1]), &meta)
	if err != nil {
		log.Fatal(err)
	}

	var htmlOutput bytes.Buffer
	err = goldmark.Convert([]byte(parts[2]), &htmlOutput)
	if err != nil {
		log.Fatal(err)
	}

	baseName := strings.TrimSuffix(filename, ".md")
	link := baseName + "/"
	if meta.Slug != "" {
		link = meta.Slug + "/"
	}

	return PostData{
		Title:           meta.Title,
		Date:            meta.Date,
		Tags:            meta.Tags,
		Link:            link,
		Content:         template.HTML(htmlOutput.String()),
		GitHubUrl:       meta.GitHubUrl,
		AdventOfCodeUrl: meta.AdventOfCodeUrl,
		LiveDemo:        meta.LiveDemo,
	}, meta.Type
}

func generateHTML(tmpl *template.Template, data interface{}, folderName string, fileName string) {
	outputDir := filepath.Join("public", folderName)

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create folder %s: %v", outputDir, err)
	}

	outputPath := filepath.Join(outputDir, fileName)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", outputPath, err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		log.Fatalf("Failed to execute template for %s: %v", outputPath, err)
	}
}

func copyStaticAssets() {
	staticFiles, err := os.ReadDir("static")
	if err != nil {
		return
	}

	for _, f := range staticFiles {
		if !f.IsDir() {
			fileBytes, _ := os.ReadFile(filepath.Join("static", f.Name()))
			os.WriteFile(filepath.Join("public", f.Name()), fileBytes, 0644)
		}
	}
}
