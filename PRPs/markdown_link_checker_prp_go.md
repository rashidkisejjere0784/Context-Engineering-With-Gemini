# Product Requirements Prompt (PRP) - Go Edition
## 1. Overview
- **Feature Name:** Markdown Link Checker CLI (Go)

- **Objective:** Build a command-line tool in Go that recursively scans a directory for Markdown files and validates all hyperlinks within them.

- **Why:** To help maintain the integrity of documentation and notes by automatically finding and reporting broken links, which is a common problem in large projects.

## 2. Success Criteria
- [ ] The code runs without errors.

- [ ] All new unit tests pass.

- [ ] The CLI tool accepts a directory path as an argument.

- [ ] The tool correctly finds all `.md` files in the given directory and its subdirectories.

- [ ] The tool extracts all local and web links from the Markdown files.

- [ ] The tool checks if local file links point to existing files.

- [ ] The tool makes HTTP requests to check if web links are live (return a 200-level status code).

- [ ] The tool prints a clear, final report of all broken links found, including the file and line number where the broken link was found.

- [ ] The code adheres to Go best practices and project standards defined in `GEMINI.md`.

## 3. Context & Resources
### 📚 External Documentation:
- **Resource:** Go `flag` package
    - **Purpose:** To build the command-line interface for accepting the directory path.

- **Resource:** Go `filepath.Walk` function
    - **Purpose:** To recursively walk through the directory structure and find all files.

- **Resource:** Go `regexp` package
    - **Purpose:** To parse the Markdown content and reliably extract link URLs. A good regex for Markdown links is `\[.*?\]\((.*?)\)`.

- **Resource:** Go `net/http` package
    - **Purpose:** For making HTTP `HEAD` requests to efficiently check web links without downloading the full content.

### 💻 Internal Codebase Patterns:
- **File:** N/A
    - **Reason:** This will be the first Go feature in the project. We will establish the patterns now.

### ⚠️ Known Pitfalls:
- Network requests can be slow. The tool should use HTTP `HEAD` requests instead of `GET`. Concurrency with Goroutines could be a future enhancement to speed up checking multiple web links.

- A website might be temporarily down. The tool should handle connection timeouts and other network errors gracefully.

- Parsing with regex can be tricky. Ensure the regex correctly handles different link formats.

## 4. Implementation Blueprint
### Proposed File Structure:
```
   src/
   └── linkchecker/
   ├── main.go               (new, handles CLI and orchestration)
   ├── parser.go             (new, handles finding files and links)
   └── validator.go          (new, handles checking if links are valid)
   tests/
   └── linkchecker_test.go     (new)
```

### Task Breakdown:
**Task 1: File & Link Parsing (`src/linkchecker/parser.go`)**

- Create a function `FindMarkdownFiles(root string) ([]string, error)`.
    - Use `filepath.Walk` to find all files ending in `.md`.

- Create a struct `LinkInfo{URL string; Line int}`.

- Create a function `ExtractLinks(filePath string) ([]LinkInfo, error)`.
    - Read the file content line by line.
    - Use `regexp.Compile` and `FindAllStringSubmatch` to extract all link URLs.
    - Return a slice of `LinkInfo` structs.

**Task 2: Link Validation (`src/linkchecker/validator.go`)**

- Create a function `ValidateLink(link LinkInfo, baseFilePath string) (bool, error)`.
    - Check if the link URL starts with `http` or `https`.
    - If it's a web URL, use `http.NewRequest("HEAD", ...)` and `http.DefaultClient.Do()` to check its status. Return `true` for 2xx status codes. Handle errors.
    - If it's a local file path, use `filepath.Dir` and `filepath.Join` to get the absolute path, then use `os.Stat` to check if the file exists.

**Task 3: Main CLI Logic (`src/linkchecker/main.go`)**

- Create a `main()` function.

- Use the `flag` package to define and parse one required argument for the target directory.

- Call `parser.FindMarkdownFiles()` to get a list of files.

- Create a struct `BrokenLink{File string; Link LinkInfo; Reason string}`.

- Loop through each file:
    - Call `parser.ExtractLinks()`.
    - Loop through each extracted link:
      - Call `validator.ValidateLink()`.
      - If a link is broken, add it to a slice of `BrokenLink` structs.

- After checking all files, print a summary report of all broken links found.

## 5. Validation Plan
### Unit Tests (`tests/linkchecker_test.go`):
- `TestFindMarkdownFiles():` Test that it correctly finds .md files in a mock directory structure.

- `TestExtractLinks():` Test with sample markdown text to ensure it extracts all links and their line numbers correctly.

- `TestValidateLocalLink():` Test both existing and non-existing local file links.

- `TestValidateWebLink():` Use Go's httptest package to create a mock HTTP server.
    - Test a link that returns `200 OK`.
    - Test a link that returns `404 Not Found`.

### Manual Test Command:
_First, create a test directory with a few markdown files, some with good links and some with bad links._

```
# Build the executable
go build -o linkchecker ./src/linkchecker/

# Run the tool
./linkchecker --directory=./my_test_docs
```
**Expected Output:**
```
Scanning for broken links in ./my_test_docs...

Found 2 broken links:
--------------------------------------------------
File: my_test_docs/file1.md (Line: 15)
Link: http://a-definitely-broken-link-12345.com
Reason: client error: 404 Not Found

File: my_test_docs/another/file2.md (Line: 8)
Link: ../non_existent_file.md
Reason: file does not exist
--------------------------------------------------
Scan complete.
```