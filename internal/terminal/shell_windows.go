//go:build windows

package terminal

import (
	"os/exec"
	"syscall"
)

// killProcessGroupSig kills the process on Windows.
// Windows doesn't have POSIX process groups, so we just kill the process.
func killProcessGroupSig(cmd *exec.Cmd, sig syscall.Signal) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	cmd.Process.Kill()
}
