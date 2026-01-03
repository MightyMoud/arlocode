package tools

import (
	"os"
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

func TestAppendLines(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	args := AppendLinesArgs{
		Path:  tmpfile.Name(),
		Lines: []string{"line1", "line2"},
	}

	result, err := AppendLines(args)
	if err != nil {
		t.Fatalf("AppendLines failed: %v", err)
	}

	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := "line1\nline2\n"
	if string(content) != expected {
		t.Errorf("Expected content '%s', got '%s'", expected, string(content))
	}

	if result == "" {
		t.Error("Expected result string, got empty")
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

	args := ReadFileArgs{
		Path: tmpfile.Name(),
	}

	result, err := ReadFileFn(args)
	if err != nil {
		t.Fatalf("ReadFileFn failed: %v", err)
	}

	if result != content {
		t.Errorf("Expected content '%s', got '%s'", content, result)
	}
}
