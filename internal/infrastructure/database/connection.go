package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type SimilarParagraph struct {
	ID       int64
	Tokens   int
	Content  string
	Distance float64
}

type Connection struct {
	Conn *pgx.Conn
}

func NewConnection(host string, port int, user string, password string, dbname string) (*Connection, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &Connection{Conn: conn}, nil
}

func (c *Connection) Close(ctx context.Context) {
	c.Conn.Close(ctx)
}
