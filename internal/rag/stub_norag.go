//go:build norag

package rag

import (
	"context"
	"time"

	"clawbench/internal/model"
)

// Stub types — no DuckDB dependency when norag build tag is set.

type Store struct{}

type Chunk struct {
	ID                 int64     `json:"id"`
	SessionID          string    `json:"session_id"`
	MessageID          int64     `json:"message_id"`
	ChunkText          string    `json:"chunk_text"`
	ChunkTextSegmented string    `json:"chunk_text_segmented"`
	ChunkIndex         int       `json:"chunk_index"`
	TokenCount         int       `json:"token_count"`
	Embedding          []float64 `json:"embedding"`
	HasEmbedding       bool      `json:"has_embedding"`
	ProjectPath        string    `json:"project_path"`
	Backend            string    `json:"backend"`
	Role               string    `json:"role"`
	CreatedAt          time.Time `json:"created_at"`
}

type SearchHit struct {
	ChunkID      int64     `json:"chunk_id"`
	ChunkText    string    `json:"chunk_text"`
	Score        float64   `json:"score"`
	SessionID    string    `json:"session_id"`
	SessionTitle string    `json:"session_title"`
	MessageID    int64     `json:"message_id"`
	Role         string    `json:"role"`
	ProjectPath  string    `json:"project_path"`
	Backend      string    `json:"backend"`
	CreatedAt    time.Time `json:"created_at"`
}

type SearchMode string

const (
	SearchModeHybrid SearchMode = "hybrid"
	SearchModeVector SearchMode = "vector"
	SearchModeFTS    SearchMode = "fts"
)

type SearchParams struct {
	Query            string `json:"q"`
	Limit            int    `json:"limit"`
	ProjectPath      string `json:"project"`
	Backend          string `json:"backend"`
	Role             string `json:"role"`
	SessionID        string `json:"session_id"`
	ExcludeSessionID string `json:"exclude_session_id"`
	FromTime         string `json:"from"`
	ToTime           string `json:"to"`
}

type SearchResult struct {
	Results []SearchHit `json:"results"`
	Total   int         `json:"total"`
	Mode    SearchMode  `json:"mode"`
}

type EmbeddingClient struct{}

type Indexer struct{}

type CleanupWorker struct{}

type TextChunk struct {
	Text       string
	TokenCount int
	Index      int
}

type PendingChunk struct {
	ID        int64
	ChunkText string
}

// Global variables — all nil in norag mode.

var (
	GlobalStore         *Store
	GlobalIndexer       *Indexer
	GlobalEmbedder      *EmbeddingClient
	GlobalCleanupWorker *CleanupWorker
)

// Stub functions — all no-ops or safe defaults.

func Init(cfg model.RAGConfig) error { return nil }

func InitStore(cfg model.RAGConfig) (*Store, error) { return nil, nil }

func NewStore(dbPath string, duckdbOpts map[string]string) (*Store, error) { return nil, nil }

func StartIndexer(cfg model.RAGConfig) {}

func StartCleanupWorker(cfg model.RAGConfig) {}

func Shutdown() {}

func EmbedderHealthy() bool { return false }

func SetEmbedderHealthy(healthy bool) {}

func RAGSearch(ctx context.Context, store *Store, embedder *EmbeddingClient, params SearchParams, defaultLimit int, searchPoolSize int) (*SearchResult, error) {
	return &SearchResult{Results: []SearchHit{}}, nil
}

func InitSegmenter() error { return nil }

func SegmentText(text string) string { return text }

func ExtractTextFromContent(content, role string) string { return "" }

func ChunkText(text string, chunkSize, chunkOverlap int) []TextChunk { return nil }

func (s *Store) Close() error { return nil }

func (s *Store) ChunkCount() (int, error) { return 0, nil }
