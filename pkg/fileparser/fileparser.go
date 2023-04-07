package fileparser

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/unidoc/unioffice/document"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func getFileExtension(filename string) string {
	return strings.ToLower(filename[strings.LastIndex(filename, "."):])
}
func isTitle(line string) bool {
	titleRegex := `^[A-Z][a-zA-Z]*( [A-Z][a-zA-Z]+)*\.?$`
	match, err := regexp.MatchString(titleRegex, line)
	if err != nil {
		return false
	}
	return match
}

func parseMarkdown(filepath string) ([]string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	md := goldmark.New(goldmark.WithExtensions())
	parsed := md.Parser().Parse(text.NewReader(data))

	var paragraphs []string
	ast.Walk(parsed, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if para, ok := n.(*ast.Paragraph); ok {
				segment := para.Lines().At(0)
				text := strings.TrimSpace(string(data[segment.Start:segment.Stop]))
				if text != "" {
					paragraphs = append(paragraphs, text)
				}
			}
		}
		return ast.WalkContinue, nil
	})

	return paragraphs, nil
}

func parseDoc(filepath string) ([]string, error) {
	doc, err := document.Open(filepath)
	if err != nil {
		return nil, err
	}

	var paragraphs []string
	for _, para := range doc.Paragraphs() {
		paraContent := ""
		for _, run := range para.Runs() {
			paraContent += run.Text()
		}
		paraContent = strings.TrimSpace(paraContent)
		if paraContent != "" && !isTitle(paraContent) {
			paragraphs = append(paragraphs, paraContent)
		}
	}
	return paragraphs, nil
}

func parseTxt(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paragraphs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !isTitle(line) {
			paragraphs = append(paragraphs, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return paragraphs, nil
}

// ParseFile takes a file path as an argument and returns an array of parsed paragraphs.
func ParseFile(filepath string) ([]string, error) {
	switch ext := getFileExtension(filepath); ext {
	case ".md":
		return parseMarkdown(filepath)
	case ".doc", ".docx":
		return parseDoc(filepath)
	case ".txt":
		return parseTxt(filepath)
	default:
		return nil, errors.New("unsupported file type")
	}
}
