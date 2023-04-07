# Text Parser and Embedding API, GPT response

This project creates an api to query with GPT on your files, TODO: pdf integration
This project also provides a text parser that can parse different file types, such as .md, .doc, .docx, and .txt, into paragraphs.
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

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
