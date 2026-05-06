package ai

// buildQoderStreamArgs constructs the CLI arguments for Qoder streaming
func buildQoderStreamArgs(req ChatRequest) []string {
	args := []string{
		"--print",
		"--output-format", "stream-json",
	}

	if req.Resume {
		args = append(args, "--resume", req.SessionID)
	} else if req.SessionID != "" {
		args = append(args, "--session-id", req.SessionID)
	}

	if req.WorkDir != "" {
		args = append(args, "--cwd", req.WorkDir)
	}

	args = append(args, "--dangerously-skip-permissions")

	// Disable built-in scheduling/timer tools to force use of ClawBench's
	// <schedule-proposal> mechanism instead of native CronCreate/CronDelete/CronList.
	args = append(args, "--disallowed-tools", "CronCreate,CronDelete,CronList")

	if req.SystemPrompt != "" {
		args = append(args, "--system-prompt", req.SystemPrompt)
	}

	// Pass model name if per-request override is set
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}

	return args
}
