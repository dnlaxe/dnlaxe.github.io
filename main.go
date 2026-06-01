package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"go.yaml.in/yaml/v3"
)

type PostData struct {
	Title   string
	Date    string
	Link    string
	Content template.HTML
}

func main() {
	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		log.Fatal(err)
	}

	indexTmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir("content")
	if err != nil {
		log.Fatal(err)
	}

	var allPosts []PostData

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}

		fmt.Printf("Processing: %s...\n", file.Name())

		inputPath := filepath.Join("content", file.Name())
		fileBytes, err := os.ReadFile(inputPath)
		if err != nil {
			log.Fatal(err)
		}

		parts := strings.SplitN(string(fileBytes), "---", 3)
		if len(parts) < 3 {
			log.Fatalf("Error: File %s is missing frontmatter", file.Name())
		}

		var post PostData
		err = yaml.Unmarshal([]byte(parts[1]), &post)
		if err != nil {
			log.Fatal(err)
		}

		var htmlOutput bytes.Buffer
		err = goldmark.Convert([]byte(parts[2]), &htmlOutput)
		if err != nil {
			log.Fatal(err)
		}

		post.Content = template.HTML(htmlOutput.String())

		baseName := strings.TrimSuffix(file.Name(), ".md")
		outputName := baseName + ".html"
		post.Link = outputName
		outputPath := filepath.Join("public", outputName)

		outputFile, err := os.Create(outputPath)
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(outputFile, post)
		if err != nil {
			log.Fatal(err)
		}
		outputFile.Close()

		allPosts = append(allPosts, post)
	}

	fmt.Println("Generating homepage...")
	indexFile, err := os.Create("public/index.html")
	if err != nil {
		log.Fatal(err)
	}

	err = indexTmpl.Execute(indexFile, allPosts)
	if err != nil {
		log.Fatal(err)
	}
	indexFile.Close()

	fmt.Println("Site fully generated.")
}
