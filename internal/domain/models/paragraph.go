package models

import "strings"

type Paragraph struct {
	ID         int64     `json:"id"`
	Content    string    `json:"content"`
	TokenCount int       `json:"token_count"`
	Embedding  []float32 `json:"embedding"`
}

func NewParagraph(content string, embedding []float32) Paragraph {
	return Paragraph{
		Content:    content,
		TokenCount: len(strings.Fields(content)),
		Embedding:  embedding,
	}
}
