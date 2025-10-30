package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// LinkInfo holds the URL and the line number where a link was found.
type LinkInfo struct {
	URL  string
	Line int
}

// FindMarkdownFiles recursively finds all markdown files in a given directory.
func FindMarkdownFiles(root string) ([]string, error) {
	var markdownFiles []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".md" {
			markdownFiles = append(markdownFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", root, err)
	}
	return markdownFiles, nil
}

// ExtractLinks reads a file and extracts all markdown-style links.
func ExtractLinks(filePath string) ([]LinkInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", filePath, err)
	}
	defer file.Close()

	var links []LinkInfo
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	// Regex to find markdown links: [text](url)
	re := regexp.MustCompile(`\[.*?\]\((.*?)\)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) > 1 {
				links = append(links, LinkInfo{URL: match[1], Line: lineNumber})
			}
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file %s: %w", filePath, err)
	}

	return links, nil
}
