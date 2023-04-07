package fileparser

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

func GetFileExtension(filename string) string {
	return strings.ToLower(filename[strings.LastIndex(filename, "."):])
}
func IsTitle(line string) bool {
	titleRegex := `^[A-Z][a-zA-Z]*( [A-Z][a-zA-Z]+)*\.?$`
	match, err := regexp.MatchString(titleRegex, line)
	if err != nil {
		return false
	}
	return match
}

// ReadFileInChunks reads a file and separates it into chunks of the specified size,
// and returns an array of strings containing the chunks.
// chunkSize is the maximum number of words per chunk, and overlap is the number of chunks
// to overlap between consecutive chunks.
func ParseTxtInChunks(filePath string, chunkSize int, overlap int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunks []string
	var words []string
	var currentWord strings.Builder

	reader := bufio.NewReader(file)

	for {
		b, err := reader.ReadByte()

		if err == nil || err == io.EOF {
			if b == ' ' || b == '\n' {
				word := currentWord.String()
				if len(word) > 0 {
					words = append(words, word)
				}
				currentWord.Reset()
			} else {
				currentWord.WriteByte(b)
			}

			// If EOF reached, break the loop after processing the last word
			if err == io.EOF {
				break
			}
		} else if err != nil {
			return nil, err
		}

		// Check if we have enough words for a chunk
		if len(words) >= chunkSize {
			// Replace newline characters with spaces and replace double quotes with two single quotes
			chunk := strings.Replace(strings.Join(words[:chunkSize], " "), "\n", " ", -1)
			chunk = strings.Replace(chunk, "\"", "''", -1)
			chunks = append(chunks, chunk)

			// Remove the first "chunkSize - overlap" words
			words = words[chunkSize-overlap:]
		}
	}

	// Add the remaining words as the last chunk
	if len(words) > 0 {
		chunk := strings.Replace(strings.Join(words, " "), "\n", " ", -1)
		chunk = strings.Replace(chunk, "\"", "''", -1)
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}
