package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/Abraxas-365/ia-search/internal/domain/models"
	"github.com/Abraxas-365/ia-search/internal/domain/ports"
	"github.com/Abraxas-365/ia-search/pkg/fileparser"
	"github.com/Abraxas-365/ia-search/pkg/openaiapi"
)

type Application interface {
	SaveParagraph(ctx context.Context, paragraph string) error
	GetGptResposeWithContext(ctx context.Context, question string) (string, error)
	ParseFile(ctx context.Context, path string) error
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

func (a *application) GetGptResposeWithContext(ctx context.Context, question string) (string, error) {
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
		if tokens >= 1500 {
			break
		}
		context = context + result.Content + "\n"
		tokens = tokens + result.TokenCount
	}
	prompt := fmt.Sprintf(`
You are a very enthusiastic bibliotecarian who loves to help people! 
Given the following sections, answer the question using only that information. 
If you are unsure and the answer is not explicitly written in the documentation, 
say "Sorry, I don't the aswer of that.

		Context sextions: %s, 

		Question: %s
		
		`, context, question)

	completition, err := a.openApi.GetCompletion(prompt, 1500, 0.5)
	if err != nil {
		return "", err
	}

	return completition, nil
}

func (a *application) ParseFile(ctx context.Context, path string) error {
	paragraphs, err := fileparser.ParseFile(path)
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
