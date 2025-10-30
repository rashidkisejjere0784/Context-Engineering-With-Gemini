package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary directory with test files
func createTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create some files
	os.WriteFile(filepath.Join(dir, "test1.md"), []byte("[link1](http://example.com)"), 0644)
	os.WriteFile(filepath.Join(dir, "test2.txt"), []byte("not a markdown file"), 0644)
	subDir := filepath.Join(dir, "subdir")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "test3.md"), []byte("[link2](./local.txt)"), 0644)
	os.WriteFile(filepath.Join(subDir, "local.txt"), []byte("I exist"), 0644)


	return dir
}

func TestFindMarkdownFiles(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	files, err := FindMarkdownFiles(dir)
	if err != nil {
		t.Fatalf("FindMarkdownFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 markdown files, but got %d", len(files))
	}
}

func TestExtractLinks(t *testing.T) {
	content := "This is a [good link](https://good.com).\nThis is a [broken link](./bad.md)."
	tmpFile, err := os.CreateTemp("", "test-extract-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Write([]byte(content))
	tmpFile.Close()

	links, err := ExtractLinks(tmpFile.Name())
	if err != nil {
		t.Fatalf("ExtractLinks failed: %v", err)
	}

	if len(links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(links))
	}
	if links[0].URL != "https://good.com" || links[0].Line != 1 {
		t.Errorf("Unexpected first link: %+v", links[0])
	}
	if links[1].URL != "./bad.md" || links[1].Line != 2 {
		t.Errorf("Unexpected second link: %+v", links[1])
	}
}

func TestValidateLocalLink(t *testing.T) {
	dir := createTestDir(t)
	defer os.RemoveAll(dir)

	// Test existing link
	baseFile := filepath.Join(dir, "subdir", "test3.md")
	existingLink := LinkInfo{URL: "./local.txt", Line: 1}
	isValid, err := ValidateLink(existingLink, baseFile)
	if !isValid || err != nil {
		t.Errorf("Expected local link to be valid, but got isValid=%v, err=%v", isValid, err)
	}

	// Test missing link
	missingLink := LinkInfo{URL: "./nonexistent.txt", Line: 2}
	isValid, err = ValidateLink(missingLink, baseFile)
	if isValid || err == nil {
		t.Errorf("Expected local link to be invalid, but got isValid=%v, err=%v", isValid, err)
	}
}

func TestValidateWebLink(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/good" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Test good link
	goodLink := LinkInfo{URL: server.URL + "/good", Line: 1}
	isValid, err := ValidateLink(goodLink, "")
	if !isValid || err != nil {
		t.Errorf("Expected web link to be valid, but got isValid=%v, err=%v", isValid, err)
	}

	// Test bad link
	badLink := LinkInfo{URL: server.URL + "/bad", Line: 2}
	isValid, err = ValidateLink(badLink, "")
	if isValid || err == nil {
		t.Errorf("Expected web link to be invalid, but got isValid=%v, err=%v", isValid, err)
	}
}
