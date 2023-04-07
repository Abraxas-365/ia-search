package ports

import (
	"context"

	"github.com/Abraxas-365/ia-search/internal/domain/models"
)

type ParagraphRepository interface {
	SaveParagraph(ctx context.Context, paragraph *models.Paragraph) error
	GetMostSimilarVectors(ctx context.Context, embedding []float32, limit int) ([]models.Paragraph, error)
	ContentExists(ctx context.Context, content string) (bool, error)
}
