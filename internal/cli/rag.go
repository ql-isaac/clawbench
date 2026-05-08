package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"clawbench/internal/model"
	"clawbench/internal/rag"
	"clawbench/internal/service"

	"gopkg.in/yaml.v3"
)

// RunRAGCommand dispatches "clawbench rag <subcommand>" CLI invocations.
func RunRAGCommand(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: clawbench rag <search|message|session> [options]\n")
		return 1
	}

	// Initialize config, database, and RAG if not already done
	if err := initRAG(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize RAG: %v\n", err)
		return 1
	}

	// Dispatch to subcommand
	switch args[0] {
	case "search":
		return runRAGSearch(args[1:])
	case "message":
		return runRAGMessage(args[1:])
	case "session":
		return runRAGSession(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", args[0])
		fmt.Fprintf(os.Stderr, "Usage: clawbench rag <search|message|session> [options]\n")
		return 1
	}
}

// initRAG initializes config, database, and RAG system for CLI usage.
func initRAG() error {
	// Skip if already initialized (e.g. in tests)
	if service.DB != nil && rag.GlobalStore != nil {
		return nil
	}

	absBinPath, _ := filepath.Abs(os.Args[0])
	model.BinDir = filepath.Dir(absBinPath)

	var cfg model.Config
	var presence map[string]bool
	configPath := filepath.Join(model.BinDir, "config", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join("config", "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = filepath.Join(model.BinDir, "config.yaml")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				configPath = "config.yaml"
			}
		}
	}

	data, err := os.ReadFile(configPath)
	if err == nil {
		var raw map[string]any
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return fmt.Errorf("parse config: %w", err)
		}
		presence = model.ParsePresenceMap(raw)
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("parse config: %w", err)
		}
	}
	model.ApplyDefaults(&cfg, presence)
	model.ConfigInstance = cfg

	if err := service.InitDB(); err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	// Initialize RAG system (always init for CLI — the search command needs it)
	if err := rag.Init(cfg.RAG); err != nil {
		return fmt.Errorf("init RAG: %w", err)
	}

	return nil
}

func runRAGSearch(args []string) int {
	fs := flagSet("search")
	query := fs.String("q", "", "Search query (required)")
	limit := fs.Int("limit", 0, "Number of results (default from config)")
	project := fs.String("project", "", "Filter by project path")
	backend := fs.String("backend", "", "Filter by backend name")
	role := fs.String("role", "", "Filter by role: user or assistant")
	sessionID := fs.String("session-id", "", "Limit results to this session")
	excludeSessionID := fs.String("exclude-session-id", "", "Exclude this session from results")
	fromTime := fs.String("from", "", "Time range start")
	toTime := fs.String("to", "", "Time range end")
	fs.Parse(args)

	if *query == "" {
		return outputError("missing required flag: -q (search query)")
	}

	if rag.GlobalStore == nil || rag.GlobalEmbedder == nil {
		return outputError("RAG is not enabled or not initialized")
	}

	params := rag.SearchParams{
		Query:            *query,
		Limit:            *limit,
		ProjectPath:      *project,
		Backend:          *backend,
		Role:             *role,
		SessionID:        *sessionID,
		ExcludeSessionID: *excludeSessionID,
		FromTime:         *fromTime,
		ToTime:           *toTime,
	}

	defaultLimit := model.ConfigInstance.RAG.SearchLimit
	if defaultLimit <= 0 {
		defaultLimit = 5
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := rag.RAGSearch(ctx, rag.GlobalStore, rag.GlobalEmbedder, params, defaultLimit)
	if err != nil {
		return outputError(fmt.Sprintf("search failed: %v", err))
	}

	if result.Results == nil {
		result.Results = []rag.SearchHit{}
	}

	outputJSON(result)
	return 0
}

func runRAGMessage(args []string) int {
	fs := flagSet("message")
	idStr := fs.String("id", "", "Message database ID (required)")
	fs.Parse(args)

	if *idStr == "" {
		// Also accept positional arg
		if fs.NArg() > 0 {
			idStr = &fs.Args()[0]
			// Can't take address of Args element directly, use a copy
			v := fs.Args()[0]
			idStr = &v
		} else {
			return outputError("missing required flag: --id (message ID)")
		}
	}

	id, err := strconv.ParseInt(*idStr, 10, 64)
	if err != nil {
		return outputError(fmt.Sprintf("invalid message ID: %v", err))
	}

	msg, err := service.GetMessageByID(id)
	if err != nil {
		return outputError(fmt.Sprintf("message not found: %v", err))
	}

	outputJSON(msg)
	return 0
}

func runRAGSession(args []string) int {
	fs := flagSet("session")
	sessionID := fs.String("id", "", "Session ID (required)")
	fs.Parse(args)

	if *sessionID == "" {
		if fs.NArg() > 0 {
			v := fs.Args()[0]
			sessionID = &v
		} else {
			return outputError("missing required flag: --id (session ID)")
		}
	}

	messages, err := service.GetMessagesBySessionID(*sessionID)
	if err != nil {
		return outputError(fmt.Sprintf("session not found: %v", err))
	}

	if messages == nil {
		messages = []model.ChatMessage{}
	}

	outputJSON(map[string]any{
		"session_id": *sessionID,
		"messages":   messages,
		"total":      len(messages),
	})
	return 0
}

// ragCLIInit initializes RAG for CLI usage in tests.
// It sets up config, DB, and RAG globals without reading config files.
func ragCLIInit(cfg model.RAGConfig) error {
	tmpDir, _ := os.MkdirTemp("", "rag-cli-test-*")
	model.BinDir = tmpDir
	model.ConfigInstance = model.Config{
		WatchDir: tmpDir,
		RAG:      cfg,
	}

	if err := service.InitDB(); err != nil {
		return err
	}

	slog.Info("rag cli: initializing RAG for test")
	return rag.Init(cfg)
}
