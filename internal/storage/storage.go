package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	GetDataFromDbByID(ctx context.Context, id uint64) (string, error)
}

type storage struct {
	conn *pgxpool.Pool
}

func (s *storage) GetDataFromDbByID(ctx context.Context, id uint64) (string, error) {
	//var name string
	//query := "SELECT name FROM students WHERE id=$1"
	//
	//ctxDb, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()
	//
	//if err := s.conn.QueryRow(ctxDb, query, id).Scan(&name); err != nil {
	//	return "", err
	//}

	return fmt.Sprintf("ID: %d\nName: John", id), nil
}

func New(conn *pgxpool.Pool) Storage {
	return &storage{
		conn: conn,
	}
}
