package sandbox

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Result struct {
	Diff    string
	RootDir string
}

func Run(cmd []string) Result {
	containerName := "safe-install-tmp"

	commandStr := join(cmd)

	// Run container
	runCmd := exec.Command("podman", "run", "--name", containerName, "-d", "ubuntu:24.04", "bash", "-c", commandStr)
	runCmd.Run()

	// Wait briefly (simple MVP)
	exec.Command("sleep", "2").Run()

	// Get diff
	diffCmd := exec.Command("podman", "diff", containerName)
	var out bytes.Buffer
	diffCmd.Stdout = &out
	diffCmd.Run()

	// Export filesystem
	rootDir := "/tmp/safe-install-root"
	exec.Command("rm", "-rf", rootDir).Run()
	exec.Command("mkdir", "-p", rootDir).Run()

	exportCmd := exec.Command("podman", "export", containerName)
	tarCmd := exec.Command("tar", "-xC", rootDir)

	pipe, _ := exportCmd.StdoutPipe()
	tarCmd.Stdin = pipe

	exportCmd.Start()
	tarCmd.Run()
	exportCmd.Wait()

	// Cleanup
	exec.Command("podman", "rm", "-f", containerName).Run()

	return Result{
		Diff:    out.String(),
		RootDir: rootDir,
	}
}

func join(cmd []string) string {
	var b bytes.Buffer
	for _, c := range cmd {
		b.WriteString(c)
		b.WriteString(" ")
	}
	return b.String()
}
