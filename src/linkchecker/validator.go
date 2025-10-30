package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ValidateLink checks if a link is valid.
// It checks local files for existence and web URLs for a 2xx status code.
func ValidateLink(link LinkInfo, baseFilePath string) (bool, error) {
	if strings.HasPrefix(link.URL, "http://") || strings.HasPrefix(link.URL, "https://") {
		// It's a web link
		return validateWebLink(link.URL)
	}
	// It's a local file link
	return validateLocalLink(link.URL, baseFilePath)
}

// validateWebLink checks if a web URL is reachable.
func validateWebLink(url string) (bool, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, fmt.Errorf("could not create request for %s: %w", url, err)
	}
	// Set a user-agent to avoid being blocked by some servers
	req.Header.Set("User-Agent", "markdown-link-checker/1.0")


	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed for %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}

	return false, fmt.Errorf("client error: %s", resp.Status)
}

// validateLocalLink checks if a local file path exists.
func validateLocalLink(linkPath string, baseFilePath string) (bool, error) {
	// If the link is absolute, check it directly. Otherwise, join it with the base path.
	if filepath.IsAbs(linkPath) {
		if _, err := os.Stat(linkPath); err == nil {
			return true, nil
		}
		return false, fmt.Errorf("file does not exist")
	}

	// Resolve the path relative to the markdown file's directory
	baseDir := filepath.Dir(baseFilePath)
	absolutePath := filepath.Join(baseDir, linkPath)

	if _, err := os.Stat(absolutePath); err == nil {
		return true, nil
	}

	return false, fmt.Errorf("file does not exist")
}
