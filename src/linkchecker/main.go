package main

import (
	"flag"
	"fmt"
	"os"
)

// BrokenLink holds information about a link that failed validation.
type BrokenLink struct {
	File   string
	Link   LinkInfo
	Reason string
}

func main() {
	// 1. Parse command-line arguments
	directory := flag.String("directory", ".", "The directory to scan for markdown files.")
	flag.Parse()

	fmt.Printf("Scanning for broken links in %s...\n\n", *directory)

	// 2. Find all markdown files
	markdownFiles, err := FindMarkdownFiles(*directory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding markdown files: %v\n", err)
		os.Exit(1)
	}

	if len(markdownFiles) == 0 {
		fmt.Println("No markdown files found.")
		return
	}

	// 3. Check links in each file
	var brokenLinks []BrokenLink
	for _, file := range markdownFiles {
		links, err := ExtractLinks(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not extract links from %s: %v\n", file, err)
			continue
		}

		for _, link := range links {
			isValid, validationErr := ValidateLink(link, file)
			if !isValid {
				brokenLinks = append(brokenLinks, BrokenLink{
					File:   file,
					Link:   link,
					Reason: validationErr.Error(),
				})
			}
		}
	}

	// 4. Print the final report
	if len(brokenLinks) > 0 {
		fmt.Printf("Found %d broken links:\n", len(brokenLinks))
		fmt.Println("--------------------------------------------------")
		for _, broken := range brokenLinks {
			fmt.Printf("File: %s (Line: %d)\n", broken.File, broken.Link.Line)
			fmt.Printf("Link: %s\n", broken.Link.URL)
			fmt.Printf("Reason: %s\n\n", broken.Reason)
		}
		fmt.Println("--------------------------------------------------")
	} else {
		fmt.Println("No broken links found. Great work!")
	}

	fmt.Println("Scan complete.")
}
