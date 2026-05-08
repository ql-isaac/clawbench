package cli

import (
	"testing"

	"clawbench/internal/model"
	"clawbench/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunRAGCommand_NoArgs(t *testing.T) {
	exitCode := RunRAGCommand([]string{})
	assert.Equal(t, 1, exitCode)
}

func TestRunRAGCommand_UnknownSubcommand(t *testing.T) {
	exitCode := RunRAGCommand([]string{"foo"})
	assert.Equal(t, 1, exitCode)
}

func TestRAGSearch_MissingQuery(t *testing.T) {
	exitCode := RunRAGCommand([]string{"search"})
	assert.Equal(t, 1, exitCode)
}

func TestRAGMessage_MissingID(t *testing.T) {
	exitCode := RunRAGCommand([]string{"message"})
	assert.Equal(t, 1, exitCode)
}

func TestRAGMessage_InvalidID(t *testing.T) {
	exitCode := RunRAGCommand([]string{"message", "--id", "abc"})
	assert.Equal(t, 1, exitCode)
}

func TestRAGSession_MissingID(t *testing.T) {
	exitCode := RunRAGCommand([]string{"session"})
	assert.Equal(t, 1, exitCode)
}

func TestRAGSession_NotFound(t *testing.T) {
	setupTestEnv(t) // reuse task test helper for DB init

	exitCode := RunRAGCommand([]string{"session", "--id", "nonexistent-session-id"})
	// Session returns empty messages, not an error — so this should succeed
	assert.Equal(t, 0, exitCode)
}

func TestRAGMessage_NotFound(t *testing.T) {
	setupTestEnv(t)

	exitCode := RunRAGCommand([]string{"message", "--id", "999999"})
	assert.Equal(t, 1, exitCode) // message not found
}

// setupTestEnv is shared with task_test.go — it initializes a temp DB.
// We need it here for message/session tests that don't need RAG initialized.
func setupRAGTestEnv(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	model.BinDir = tmpDir
	model.ConfigInstance = model.Config{
		WatchDir: tmpDir,
		RAG: model.RAGConfig{
			OllamaBaseURL: "http://localhost:11434",
			OllamaModel:   "bge-m3",
			SearchLimit:   5,
		},
	}

	err := service.InitDB()
	require.NoError(t, err)
}
