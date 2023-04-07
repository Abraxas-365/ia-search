package repository

import (
	"context"

	"github.com/Abraxas-365/ia-search/internal/domain/models"
	"github.com/Abraxas-365/ia-search/internal/domain/ports"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type paragraphRepository struct {
	db *pgxpool.Pool
}

func NewParagraphRepository(db *pgxpool.Pool) ports.ParagraphRepository {
	return &paragraphRepository{db: db}
}

func (r *paragraphRepository) SaveParagraph(ctx context.Context, paragraph *models.Paragraph) error {
	query := `INSERT INTO "public"."paragraph" (content, token_count, embedding) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, paragraph.Content, paragraph.TokenCount, pgvector.NewVector(paragraph.Embedding))
	return err
}

func (r *paragraphRepository) GetMostSimilarVectors(ctx context.Context, embedding []float32, limit int) ([]models.Paragraph, error) {
	query := `
	SELECT id, token_count, content
	FROM "public"."paragraph"
	ORDER BY embedding <-> $1
	LIMIT $2;
	`

	rows, err := r.db.Query(ctx, query, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var similarParagraphs []models.Paragraph

	for rows.Next() {
		var sp models.Paragraph

		err := rows.Scan(&sp.ID, &sp.TokenCount, &sp.Content)
		if err != nil {
			return nil, err
		}

		similarParagraphs = append(similarParagraphs, sp)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return similarParagraphs, nil
}

func (r *paragraphRepository) ContentExists(ctx context.Context, content string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM "public"."paragraph" WHERE content = $1)`
	err := r.db.QueryRow(ctx, query, content).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
