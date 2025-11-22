package monitor

import (
	"os/exec"
	"syscall"
	"time"
)

const ProcessKillDelay = 500 * time.Millisecond

func setProcessGroup(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killProcessTree(pid int) error {
	_ = syscall.Kill(-pid, syscall.SIGTERM)
	time.Sleep(ProcessKillDelay)
	_ = syscall.Kill(-pid, syscall.SIGKILL)
	return nil
}
