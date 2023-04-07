package fileparser

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
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
	wordChan := make(chan string)
	done := make(chan struct{})
	errChan := make(chan error)

	go func() {
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
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case word, ok := <-wordChan:
				if !ok {
					done <- struct{}{}
					return
				}
				words = append(words, word)

				if len(words) >= chunkSize {
					chunk := strings.Replace(strings.Join(words[:chunkSize], " "), "\n", " ", -1)
					chunk = strings.Replace(chunk, "\"", "''", -1)
					chunks = append(chunks, chunk)

					words = words[chunkSize-overlap:]
				}
			case err := <-errChan:
				done <- struct{}{}
				errChan <- err
				return
			}
		}
	}()

	<-done
	wg.Wait()

	if len(words) > 0 {
		chunk := strings.Replace(strings.Join(words, " "), "\n", " ", -1)
		chunk = strings.Replace(chunk, "\"", "''", -1)
		chunks = append(chunks, chunk)
	}

	select {
	case err := <-errChan:
		return nil, err
	default:
		return chunks, nil
	}
}
