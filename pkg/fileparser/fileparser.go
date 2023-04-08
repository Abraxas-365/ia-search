package fileparser

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

// StopWords is a map containing common stop words
var stopWords = map[string]bool{
	"a":         true,
	"about":     true,
	"above":     true,
	"after":     true,
	"again":     true,
	"against":   true,
	"all":       true,
	"am":        true,
	"an":        true,
	"and":       true,
	"any":       true,
	"are":       true,
	"as":        true,
	"at":        true,
	"be":        true,
	"because":   true,
	"been":      true,
	"before":    true,
	"being":     true,
	"below":     true,
	"between":   true,
	"both":      true,
	"but":       true,
	"by":        true,
	"can":       true,
	"did":       true,
	"do":        true,
	"does":      true,
	"doing":     true,
	"don't":     true,
	"down":      true,
	"during":    true,
	"each":      true,
	"few":       true,
	"for":       true,
	"from":      true,
	"further":   true,
	"had":       true,
	"has":       true,
	"have":      true,
	"having":    true,
	"he":        true,
	"he'd":      true,
	"he'll":     true,
	"he's":      true,
	"her":       true,
	"here":      true,
	"here's":    true,
	"hers":      true,
	"herself":   true,
	"him":       true,
	"himself":   true,
	"his":       true,
	"how":       true,
	"how's":     true,
	"i":         true,
	"i'd":       true,
	"i'll":      true,
	"i'm":       true,
	"i've":      true,
	"if":        true,
	"in":        true,
	"into":      true,
	"is":        true,
	"it":        true,
	"it's":      true,
	"its":       true,
	"itself":    true,
	"let's":     true,
	"me":        true,
	"more":      true,
	"most":      true,
	"my":        true,
	"myself":    true,
	"nor":       true,
	"of":        true,
	"on":        true,
	"once":      true,
	"only":      true,
	"or":        true,
	"other":     true,
	"ought":     true,
	"our":       true,
	"ours":      true,
	"ourselves": true,
	"out":       true,
	"over":      true,
	"own":       true,
	"same":      true,
	"she":       true,
	"she'd":     true,
	"she'll":    true,
	"she's":     true,
	"should":    true,
	"so":        true,
	"some":      true,
	"such":      true,
	"than":      true,
	"that":      true,
	"that's":    true,
	"the":       true,
	"their":     true,
	"theirs":    true,
}

// GetFileExtension returns the file extension in lowercase
func GetFileExtension(filename string) string {
	return strings.ToLower(filename[strings.LastIndex(filename, "."):])
}

// IsTitle checks if the given line is a title based on regex
func IsTitle(line string) bool {
	titleRegex := `^[A-Z][a-zA-Z]*( [A-Z][a-zA-Z]+)*\.?$`
	match, err := regexp.MatchString(titleRegex, line)
	if err != nil {
		return false
	}
	return match
}

// CleanWord removes all non-alphanumeric characters from a word
func CleanWord(word string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(word, "")
}

// PreprocessWords cleans and filters stop words from the given words
func PreprocessWords(words []string) []string {
	var preprocessedWords []string

	for _, word := range words {
		cleanedWord := CleanWord(word)
		lowercasedWord := strings.ToLower(cleanedWord)

		if len(lowercasedWord) > 0 && !stopWords[lowercasedWord] {
			preprocessedWords = append(preprocessedWords, lowercasedWord)
		}
	}

	return preprocessedWords
}

// ReadWords reads words from a file and sends them to a channel
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

// ParseWordsInChunks creates chunks of text and sends them to a channel
func ParseWordsInChunks(wordChan <-chan string, chunkSize int, overlap int, chunksChan chan<- string, errChan chan<- error) {
	var words []string
	wordCount := 0

	for word := range wordChan {
		words = append(words, word)
		wordCount++

		if wordCount >= chunkSize {
			// Find the position of the last period in the chunk
			lastPeriodIndex := -1
			for i := len(words) - 1; i >= 0; i-- {
				if strings.Contains(words[i], ".") {
					lastPeriodIndex = i
					break
				}
			}

			if lastPeriodIndex != -1 {
				chunk := strings.Replace(strings.Join(words[:lastPeriodIndex+1], " "), "\n", " ", -1)
				chunk = strings.Replace(chunk, "\"", "''", -1)
				chunksChan <- chunk

				// Calculate how many words to remove from the beginning of the words array
				wordsToRemove := lastPeriodIndex + 1 - overlap
				if wordsToRemove < 0 {
					wordsToRemove = 0
				}

				words = words[wordsToRemove:]
				wordCount = len(words)
			}
		}
	}

	if len(words) > 0 {
		chunk := strings.Replace(strings.Join(words, " "), "\n", " ", -1)
		chunk = strings.Replace(chunk, "\"", "''", -1)
		chunksChan <- chunk
	}

	close(chunksChan)
}

// ParseTxtInChunks reads words from a file, creates chunks, and returns them
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
