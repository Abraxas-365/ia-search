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
func ReadWords(filePath string, wordChan chan<- string, errChan chan<- error) {
	file, err := os.Open(filePath)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var currentWord strings.Builder

	for {
		b, err := reader.ReadByte()

		if err == nil || err == io.EOF {
			if b == ' ' || b == '\n' {
				word := currentWord.String()
				if len(word) > 0 {
					wordChan <- word
				}
				currentWord.Reset()
			} else {
				currentWord.WriteByte(b)
			}

			if err == io.EOF {
				close(wordChan)
				return
			}
		} else if err != nil {
			errChan <- err
			return
		}
	}
}

func ParseWordsInChunks(wordChan <-chan string, chunkSize int, overlap int, chunksChan chan<- string, errChan chan<- error) {
	var words []string

	for word := range wordChan {
		words = append(words, word)

		if len(words) >= chunkSize {
			chunk := strings.Replace(strings.Join(words[:chunkSize], " "), "\n", " ", -1)
			chunk = strings.Replace(chunk, "\"", "''", -1)
			chunksChan <- chunk

			words = words[chunkSize-overlap:]
		}
	}

	if len(words) > 0 {
		chunk := strings.Replace(strings.Join(words, " "), "\n", " ", -1)
		chunk = strings.Replace(chunk, "\"", "''", -1)
		chunksChan <- chunk
	}

	close(chunksChan)
}

func ParseTxtInChunks(filePath string, chunkSize int, overlap int) ([]string, error) {
	wordChan := make(chan string)
	chunksChan := make(chan string)
	errChan := make(chan error)

	go ReadWords(filePath, wordChan, errChan)
	go ParseWordsInChunks(wordChan, chunkSize, overlap, chunksChan, errChan)

	var chunks []string
	done := make(chan struct{})

	go func() {
		for chunk := range chunksChan {
			chunks = append(chunks, chunk)
		}
		done <- struct{}{}
	}()

	<-done

	select {
	case err := <-errChan:
		return nil, err
	default:
		return chunks, nil
	}
}
