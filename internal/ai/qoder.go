package ai

import (
	"os/exec"
	"strings"
)

// qoderBackend is the CLIBackend instance for Qoder CLI.
var qoderBackend = &CLIBackend{
	name:           "qoder",
	defaultCommand: "qodercli",
	buildArgs:      buildQoderStreamArgs,
	newParser:      func() LineParser { return &StreamParser{} },
	filterLine:     nil, // skip empty lines only (default)
	preStart: func(cmd *exec.Cmd, req ChatRequest) {
		// Qoder CLI in --print mode with stdout piped (non-TTY) requires
		// prompt via stdin — positional prompt argument is not recognized.
		// Both new sessions and resume sessions use stdin for prompt.
		cmd.Stdin = strings.NewReader(req.Prompt)
	},
}
