package tools

import (
	"fmt"
	"os"
	"strings"
)

type AppendLinesArgs struct {
	Path  string   `json:"path" jsonschema:"The absolute or relative path to the file"`
	Lines []string `json:"lines" jsonschema:"A list of strings to append as new lines"`
}

func AppendLines(args AppendLinesArgs) (string, error) {
	f, err := os.OpenFile(args.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	content := strings.Join(args.Lines, "\n") + "\n"
	if _, err := f.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}

	return fmt.Sprintf("Successfully added %d lines to %s", len(args.Lines), args.Path), nil
}

type ReadFileArgs struct {
	Path string `json:"path" jsonschema:"The file path to read"`
}

// Then pass this wrapper function to NewButlerTool
func ReadFileFn(args ReadFileArgs) (string, error) {
	bytesData, _ := os.ReadFile(args.Path)
	return string(bytesData), nil
}

var appendLinesTool = NewButlerTool("append_lines", "Adds multiple lines to the end of a file", AppendLines)
var readFileTool = NewButlerTool("read_file", "Reads a file from the user pc", ReadFileFn)
var StdToolset = []Tool{
	appendLinesTool,
	readFileTool,
}
