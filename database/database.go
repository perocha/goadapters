package database

import (
	"context"
)

// Database represents the interface for interacting with the database.
type Database interface {
	Repository() DBRepository
}

// DBRepository represents the interface for interacting with document in the database.
type DBRepository interface {
	CreateDocument(ctx context.Context, document interface{}) error
	UpdateDocument(ctx context.Context, id string, document interface{}) error
	DeleteDocument(ctx context.Context, id string, partitionKey string) error
	GetDocument(ctx context.Context, id string, partitionKey string) (interface{}, error)
}
