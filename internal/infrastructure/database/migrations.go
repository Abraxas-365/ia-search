package database

import "context"

func (c *Connection) RunMigrations(ctx context.Context) error {
	// Create extension
	_, err := c.Pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA public;")
	if err != nil {
		return err
	}

	// Create table
	_, err = c.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS "public"."paragraph" (
			id BIGSERIAL PRIMARY KEY,
			content TEXT,
			token_count int,
			embedding vector(1536)
		);
	`)
	return err
}
