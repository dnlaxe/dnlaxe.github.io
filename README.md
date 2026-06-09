# My Custom Static Site Generator

https://dnlaxe.github.io/

A static site generator (SSG) made with Go. It converts Markdown files into a HTML/CSS website. Inspired by Hugo and Kindle.

## Features

- Zero frontend dependencies
- Markdown support
- Dynamic routing
- Client filtering
- Persistent dark mode
- Component architecture
- GitHub Pages ready

## Installation

1. Clone this repository.
2. Download the required Go modules (Goldmark and YAML parser):
   ```
   go mod tidy
   ```

## Running Locally

To build the site and start server, use the --serve flag:

```
go run main.go --serve
```

Then, open your browser to http://localhost:8080. Press CTRL+C in your terminal to stop the server.

## Building for Production

To build the HTML files without starting the server (this is what GitHub Actions uses):

```
go run main.go
```

Project Structure

```
├── content/           # Markdown files go here
├── public/            # The final compiled website
├── static/            # Static assets like styles.css, images, and favicons
├── templates/         # Go HTML templates and shared partials
├── main.go            # The core Go engine
└── README.md
```

## Tech Stack

- Engine: Go (Templates, FileSystem)
- Parsers: yuin/goldmark (Markdown), gopkg.in/yaml.v3 (Frontmatter)
- Frontend: Vanilla HTML5, CSS3, JavaScript
