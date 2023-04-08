package application

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Abraxas-365/ia-search/internal/domain/models"
	"github.com/Abraxas-365/ia-search/internal/domain/ports"
	"github.com/Abraxas-365/ia-search/pkg/fileparser"
	"github.com/Abraxas-365/ia-search/pkg/openaiapi"
)

type Application interface {
	SaveParagraph(ctx context.Context, paragraph string) error
	GetGptResposeWithContext(ctx context.Context, question string, contextTokenSize int, model string, chat bool) (string, error)
	ParseFile(ctx context.Context, path string, chucks int, overlap int) error
}

type application struct {
	repo    ports.ParagraphRepository
	openApi *openaiapi.Client
}

func NewApplication(repo ports.ParagraphRepository, openApiKey string) Application {
	return &application{
		repo:    repo,
		openApi: openaiapi.NewClient(openApiKey),
	}
}

func (a *application) SaveParagraph(ctx context.Context, content string) error {
	exists, err := a.repo.ContentExists(ctx, content)
	if err != nil {
		return err
	}

	if !exists {
		embedding, err := a.openApi.GetEmbedding(content)
		if err != nil {
			return err
		}
		paragraph := models.NewParagraph(content, embedding)
		return a.repo.SaveParagraph(ctx, &paragraph)
	}
	return nil
}

func (a *application) GetGptResposeWithContext(ctx context.Context, question string, contextTokenSize int, model string, chat bool) (string, error) {
	context := ""
	tokens := 0
	embedding, err := a.openApi.GetEmbedding(question)
	if err != nil {
		return "", err
	}

	results, err := a.repo.GetMostSimilarVectors(ctx, embedding, 5)
	if err != nil {
		return "", err
	}

	for _, result := range results {
		if tokens > contextTokenSize {
			break
		}
		context = context + result.Content + "\n"
		tokens = tokens + result.TokenCount
	}
	prompt := fmt.Sprintf(`
		Answer with the text in the context, you can use additional info that you have if it helps to complete the info in the context for the question.
		If the Question doest have any realtion to the Context or if you dont know the answer , tell I cant help you with that.

		Context sextions: %s, 

		Question: %s

		Give just the aswer

		`, context, question)

	completition := ""
	if chat {
		completition, err = a.openApi.GetChatCompletion(completition, 0.4, model)
		log.Println("error in getCompletion")
		return "", err
	} else {
		completition, err = a.openApi.GetCompletion(prompt, 1500, 0.4, model)
		if err != nil {
			log.Println("error in getCompletion")
			return "", err
		}
	}

	return strings.Trim(completition, "\n"), nil
}

func (a *application) ParseFile(ctx context.Context, path string, chucks int, overlap int) error {
	paragraphs, err := fileparser.ParseTxtInChunks(path, chucks, overlap)
	fmt.Println("Found", len(paragraphs))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	var wg sync.WaitGroup
	errorChan := make(chan error, len(paragraphs))

	for _, para := range paragraphs {
		wg.Add(1)
		go func(content string) {
			defer wg.Done()
			err := a.SaveParagraph(ctx, content)
			if err != nil {
				errorChan <- err
			}
		}(para)
	}

	wg.Wait()
	close(errorChan)

	if len(errorChan) > 0 {
		// Return the first error encountered
		return <-errorChan
	}

	return nil
}
