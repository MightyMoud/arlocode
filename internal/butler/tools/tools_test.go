package tools

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewButlerTool(t *testing.T) {
	handler := func(args struct{ Name string }) (string, error) {
		return "Hello " + args.Name, nil
	}
	tool := NewButlerTool("greet", "Greets a person", handler)

	if tool.Name != "greet" {
		t.Errorf("Expected name 'greet', got '%s'", tool.Name)
	}
	if tool.Description != "Greets a person" {
		t.Errorf("Expected description 'Greets a person', got '%s'", tool.Description)
	}
}

func TestReadFileFn(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	content := "hello world"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	args := readFileArgs{
		Path: tmpfile.Name(),
	}

	result, err := readFileFn(args)
	if err != nil {
		t.Fatalf("ReadFileFn failed: %v", err)
	}

	if result != content {
		t.Errorf("Expected content '%s', got '%s'", content, result)
	}
}

func TestReadFolderContentFn(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "testfolder")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	file1Content := "content of file 1"
	file2Content := "content of file 2"

	err = os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte(file1Content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte(file2Content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	args := readFolderContentArgs{
		FolderPath: tmpDir,
	}

	result, err := readFolderContentFn(args)
	if err != nil {
		t.Fatalf("readFolderContentFn failed: %v", err)
	}

	// Check if both files are in the result
	if !strings.Contains(result, "file1.txt") {
		t.Errorf("Expected result to contain 'file1.txt'")
	}
	if !strings.Contains(result, "file2.txt") {
		t.Errorf("Expected result to contain 'file2.txt'")
	}
	if !strings.Contains(result, file1Content) {
		t.Errorf("Expected result to contain file1 content")
	}
	if !strings.Contains(result, file2Content) {
		t.Errorf("Expected result to contain file2 content")
	}
}

func TestListFolderContents(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "testlist")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files and directories
	err = os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	args := listFolderContentsArgs{
		FolderPath: tmpDir,
	}

	result, err := listFolderContents(args)
	if err != nil {
		t.Fatalf("listFolderContents failed: %v", err)
	}

	// Check if listings are present
	if !strings.Contains(result, "File: file1.txt") {
		t.Errorf("Expected result to contain 'File: file1.txt'")
	}
	if !strings.Contains(result, "Directory: subdir") {
		t.Errorf("Expected result to contain 'Directory: subdir'")
	}
	if !strings.Contains(result, "File: subdir") && !strings.Contains(result, "file2.txt") {
		t.Errorf("Expected result to contain file2.txt reference")
	}
}

func TestSearchCode(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "testsearch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files with searchable content
	file1Content := "package main\nfunc hello() {\n\tfmt.Println(\"hello\")\n}"
	file2Content := "package test\nfunc world() {\n\tfmt.Println(\"world\")\n}"

	err = os.WriteFile(filepath.Join(tmpDir, "file1.go"), []byte(file1Content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "file2.go"), []byte(file2Content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Test searching for "hello"
	args := searchCodeArgs{
		FolderPath: tmpDir,
		Query:      "hello",
	}

	result, err := searchCode(args)
	if err != nil {
		t.Fatalf("searchCode failed: %v", err)
	}

	if !strings.Contains(result, "file1.go") {
		t.Errorf("Expected result to contain 'file1.go'")
	}
	if strings.Contains(result, "file2.go") {
		t.Errorf("Expected result to NOT contain 'file2.go'")
	}

	// Test searching for non-existent query
	args.Query = "nonexistent"
	result, err = searchCode(args)
	if err != nil {
		t.Fatalf("searchCode failed: %v", err)
	}

	if result != "No matches found." {
		t.Errorf("Expected 'No matches found.', got '%s'", result)
	}
}

func TestApplyEdit(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "edittest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	originalContent := "Line 1\nLine 2\nLine 3"
	if _, err := tmpfile.WriteString(originalContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test successful edit
	args := applyEditArgs{
		Path:    tmpfile.Name(),
		OldText: "Line 2",
		NewText: "Modified Line 2",
	}

	result, err := applyEdit(args)
	if err != nil {
		t.Fatalf("applyEdit failed: %v", err)
	}

	if result != "Edit applied successfully." {
		t.Errorf("Expected success message, got '%s'", result)
	}

	// Verify the edit was applied
	newContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expectedContent := "Line 1\nModified Line 2\nLine 3"
	if string(newContent) != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, string(newContent))
	}

	// Test old text not found
	args.OldText = "Nonexistent text"
	_, err = applyEdit(args)
	if err == nil {
		t.Error("Expected error when old_text is not found")
	}

	// Test ambiguous old text
	os.WriteFile(tmpfile.Name(), []byte("duplicate\nduplicate"), 0644)
	args.OldText = "duplicate"
	args.NewText = "replaced"
	_, err = applyEdit(args)
	if err == nil {
		t.Error("Expected error when old_text appears multiple times")
	}
}

func TestFetchURLAsMarkdown(t *testing.T) {
	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Test Page</title></head>
		<body>
			<article>
				<h1>Test Article</h1>
				<p>This is a test paragraph.</p>
			</article>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	args := fetchURLAsMarkdownArgs{
		URL: server.URL,
	}

	result, err := fetchURLAsMarkdown(args)
	if err != nil {
		t.Fatalf("fetchURLAsMarkdown failed: %v", err)
	}

	// Check if the result contains expected markdown content
	if !strings.Contains(result, "Test") {
		t.Errorf("Expected result to contain 'Test'")
	}

	// Test with bad URL
	args.URL = "http://nonexistent-domain-12345.com"
	_, err = fetchURLAsMarkdown(args)
	if err == nil {
		t.Error("Expected error when fetching invalid URL")
	}
}

func TestMakeFileFn(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "testmakefile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test creating a new file
	filePath := filepath.Join(tmpDir, "newfile.txt")
	content := "This is test content\nLine 2\nLine 3"

	args := makeFileWithContentArgs{
		Path:    filePath,
		Content: content,
	}

	result, err := makeFileWithContentFn(args)
	if err != nil {
		t.Fatalf("makeFileWithContentFn failed: %v", err)
	}

	expectedResult := "File created at " + filePath
	if result != expectedResult {
		t.Errorf("Expected result '%s', got '%s'", expectedResult, result)
	}

	// Verify the file was created with correct content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	if string(fileContent) != content {
		t.Errorf("Expected file content '%s', got '%s'", content, string(fileContent))
	}

	// Test creating a file in a subdirectory that doesn't exist yet
	subDirPath := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDirPath, 0755)
	nestedFilePath := filepath.Join(subDirPath, "nested.txt")
	nestedContent := "Nested file content"

	nestedArgs := makeFileWithContentArgs{
		Path:    nestedFilePath,
		Content: nestedContent,
	}

	_, err = makeFileWithContentFn(nestedArgs)
	if err != nil {
		t.Fatalf("makeFileWithContentFn failed for nested file: %v", err)
	}

	// Verify nested file
	nestedFileContent, err := os.ReadFile(nestedFilePath)
	if err != nil {
		t.Fatalf("Failed to read nested file: %v", err)
	}

	if string(nestedFileContent) != nestedContent {
		t.Errorf("Expected nested file content '%s', got '%s'", nestedContent, string(nestedFileContent))
	}

	// Test overwriting an existing file
	overwriteContent := "Overwritten content"
	overwriteArgs := makeFileWithContentArgs{
		Path:    filePath,
		Content: overwriteContent,
	}

	_, err = makeFileWithContentFn(overwriteArgs)
	if err != nil {
		t.Fatalf("makeFileWithContentFn failed when overwriting: %v", err)
	}

	// Verify the file was overwritten
	overwrittenContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read overwritten file: %v", err)
	}

	if string(overwrittenContent) != overwriteContent {
		t.Errorf("Expected overwritten content '%s', got '%s'", overwriteContent, string(overwrittenContent))
	}
}
