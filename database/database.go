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
	CreateDocument(ctx context.Context, partitionKey string, document interface{}) error
	UpdateDocument(ctx context.Context, partitionKey string, id string, document interface{}) error
	DeleteDocument(ctx context.Context, partitionKey string, id string) error
	GetDocument(ctx context.Context, partitionKey string, id string) (interface{}, error)
}
