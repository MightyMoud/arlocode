package tools

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	readability "codeberg.org/readeck/go-readability/v2"
	md "github.com/JohannesKaufmann/html-to-markdown"
)

type readFileArgs struct {
	Path string `json:"path" jsonschema:"The file path to read"`
}

// Then pass this wrapper function to NewButlerTool
func readFileFn(args readFileArgs) (string, error) {
	bytesData, err := os.ReadFile(args.Path)
	if err != nil {
		return "", err
	}
	return string(bytesData), nil
}

type readFolderContentArgs struct {
	FolderPath string `json:"folder_path" jsonschema:"The folder path to read all files from"`
}

func readFolderContentFn(args readFolderContentArgs) (string, error) {
	var builder strings.Builder

	err := filepath.Walk(args.FolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(args.FolderPath, path)
			builder.WriteString(fmt.Sprintf("File: %s\n", relativePath))
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			builder.WriteString(string(data) + "\n\n")
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

type listFolderContentsArgs struct {
	FolderPath string `json:"folder_path" jsonschema:"The folder path to list contents of"`
}

func listFolderContents(args listFolderContentsArgs) (string, error) {
	var builder strings.Builder

	err := filepath.Walk(args.FolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		relativePath, _ := filepath.Rel(args.FolderPath, path)
		if info.IsDir() {
			builder.WriteString(fmt.Sprintf("Directory: %s\n", relativePath))
		} else {
			builder.WriteString(fmt.Sprintf("File: %s\n", relativePath))
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

type searchCodeArgs struct {
	FolderPath string `json:"folder_path" jsonschema:"The folder path to search code in"`
	Query      string `json:"query" jsonschema:"The search query to look for in the code"`
}

func searchCode(args searchCodeArgs) (string, error) {
	var files []string

	err := filepath.Walk(args.FolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			relativePath, _ := filepath.Rel(args.FolderPath, path)
			files = append(files, relativePath)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "No matches found.", nil
	}

	resultChan := make(chan string)
	var wg sync.WaitGroup

	for _, relPath := range files {
		wg.Add(1)
		go func(relPath string) {
			defer wg.Done()
			fullPath := filepath.Join(args.FolderPath, relPath)
			data, err := os.ReadFile(fullPath)
			if err != nil {
				// skip files that can't be read
				return
			}
			lines := strings.Split(string(data), "\n")
			for i, line := range lines {
				if strings.Contains(line, args.Query) {
					resultChan <- fmt.Sprintf("Found in file: %s, line: %d\n", relPath, i+1)
				}
			}
		}(relPath)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var builder strings.Builder
	for res := range resultChan {
		builder.WriteString(res)
	}

	result := builder.String()
	if result == "" {
		return "No matches found.", nil
	}
	return result, nil
}

type applyEditArgs struct {
	Path    string `json:"path" jsonschema:"The absolute or relative path to the file"`
	OldText string `json:"old_text" jsonschema:"The exact text block to find"`
	NewText string `json:"new_text" jsonschema:"The block text to replace it with"`
}

func applyEdit(req applyEditArgs) (string, error) {
	// 1. Read the file
	content, err := os.ReadFile(req.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	fileString := string(content)

	// 2. Safety Check: Does the old text exist?
	count := strings.Count(fileString, req.OldText)
	if count == 0 {
		return "", fmt.Errorf("could not find the 'old_text' block in %s. Ensure formatting and whitespace match exactly", req.Path)
	}
	if count > 1 {
		return "", fmt.Errorf("the 'old_text' block is ambiguous (found %d occurrences). Please provide more context", count)
	}

	// 3. Perform the replacement
	newContent := strings.Replace(fileString, req.OldText, req.NewText, 1)

	// 4. Write back to disk
	err = os.WriteFile(req.Path, []byte(newContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	return "Edit applied successfully.", nil
}

// Best suited for fetching docs and github readmes -> might need another general web scraper later
type fetchURLAsMarkdownArgs struct {
	URL string `json:"url" jsonschema:"The URL of the webpage to fetch and convert to markdown must include the protocol, e.g., https://example.com"`
}

func fetchURLAsMarkdown(args fetchURLAsMarkdownArgs) (string, error) {
	resp, err := http.Get(args.URL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	parsedURL, _ := url.Parse(args.URL)
	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse article: %w", err)
	}

	var HTMLContent bytes.Buffer
	article.RenderHTML(&HTMLContent)

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(HTMLContent.String())
	if err != nil {
		return "", fmt.Errorf("failed to convert to markdown: %w", err)
	}

	finalOutput := fmt.Sprintf("# %s\n\n%s", article.Title(), markdown)

	return finalOutput, nil

}

var readFileTool = NewButlerTool("read_file", "Reads a file from the user pc - do not use this to read content from a URL", readFileFn)
var readFolderTool = NewButlerTool("read_folder", "Reads all files from a folder on the user pc", readFolderContentFn)
var listFolderContentsTool = NewButlerTool("list_folder_contents", "Lists all files and directories in a folder on the user pc - this tool will only list the files and won't read them", listFolderContents)
var searchCodeTool = NewButlerTool("search_code", "Searches for a query string in all code files within a specified folder", searchCode)
var applyEditTool = NewButlerTool("apply_edit", "Applies a code edit by replacing old text with new text in a specified file", applyEdit)
var fetchURLAsMarkdownTool = NewButlerTool("fetch_url_as_markdown", "Fetches a webpage from a URL and converts its content to markdown format for the model to read and understand. This works best with github readmes and documentation pages", fetchURLAsMarkdown)

var StdToolset = []Tool{
	readFileTool,
	readFolderTool,
	listFolderContentsTool,
	searchCodeTool,
	applyEditTool,
	fetchURLAsMarkdownTool,
}
