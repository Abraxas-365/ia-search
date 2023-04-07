# Text Parser and Embedding API, GPT response

This project creates an api to query with GPT on your files, TODO: pdf integration
This project also provides a text parser, for parsing into chunks to be embbeded.
It also includes an API to interact with OpenAI's GPT-3 and text embeddings services.

## Features

- Query you files using embedings and answer your questions using gpt
- Parse various file formats into paragraphs
- Connect to OpenAI API for embeddings and completions
- Store paragraph embeddings in a PostgreSQL database using the pgvector extension
- Find most similar paragraphs based on embeddings
- Golang package structure with internal and pkg directories

### Usage

1. Put your OPEN IA KEY in the docker-compose.yml
2. In main.go put the file you want to use
3. At the moment we are usinig text-davinci-003 that only accept 2048 tokens, so we neet to tweek the config to match that pre requists

### Start Project

```bash
docker-compose build
docker-compose up -d
```

## Example API Call

You can make an API call to get completions using a REST client or a tool like `curl`. Here's an example API call to a locally running server:

- Endpoint http://localhost:3000/api/completition
- body: {"query": "string"}

```bash
curl -X POST http://localhost:3000/api/completition \
  -H "Content-Type: application/json" \
  -d '{"query": "what is the article about?"}'
```

## How its working (parser)

### Function: ParseTxtInChunks

`ParseTxtInChunks` function reads a text file located at the given filePath and splits its content into chunks of chunkSize words with an overlap number of words between each chunk. It returns a slice of strings, where each string represents a chunk of text.

## How its working (application layer)

### Function: SaveParagraph

The `SaveParagraph` method first checks if the content already exists in the database by calling the `ContentExists` method from the `repo`. If the content exists, it returns immediately, not saving the paragraph again.

If the content does not exist, the method retrieves the paragraph's embedding from the OpenAI API using the `GetEmbedding` method. Then, it creates a new `models.Paragraph` instance with the content and its embedding, and saves the paragraph in the repository using the `SaveParagraph` method.

```go
func (a *application) SaveParagraph(ctx context.Context, content string) error {
	// Check if the content already exists in the database
	exists, err := a.repo.ContentExists(ctx, content)
	if err != nil {
		return err
	}

	// If the content doesn't exist, save the paragraph with its embedding
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
```

### Function: GetGptResposeWithContext

The `GetGptResposeWithContext` method first retrieves the embedding for the given question using the OpenAI API. It then finds the most similar paragraphs to the question by calling the `GetMostSimilarVectors` method from the `repo`.

It builds a context string by concatenating the content of the most similar paragraphs, limiting the total tokens to 1500. The context string and the question are then formatted into a GPT-3 prompt. Finally, it calls the OpenAI API to get a completion based on the generated prompt, using the `GetCompletion` method.

```go
func (a *application) GetGptResposeWithContext(ctx context.Context, question string, model string) (string, error) {
	// Get the embedding for the question
	embedding, err := a.openApi.GetEmbedding(question)
	if err != nil {
		return "", err
	}

	// Find the most similar paragraphs
	results, err := a.repo.GetMostSimilarVectors(ctx, embedding, 5)
	if err != nil {
		return "", err
	}

	// Build the context string
	context := ""
	tokens := 0
	for _, result := range results {
		if tokens >= 2000 {
			break
		}
		context = context + result.Content + "\n"
		tokens = tokens + result.TokenCount
	}

	// Format the GPT-3 prompt
	prompt := fmt.Sprintf(`
		// ...
		Context sextions: %s,
		Question: %s
		`, context, question)

	// Get the completion from the OpenAI API
	completion, err := a.openApi.GetCompletion(prompt, 1500, 0.5, model)
	if err != nil {
		return "", err
	}

	return completion, nil
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
