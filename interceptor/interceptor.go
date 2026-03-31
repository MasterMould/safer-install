package interceptor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"safe-install/sandbox"
)

type Action struct {
	Intercept bool
	Reason    string
	Command   []string
}

func Check(args []string) Action {
	cmd := strings.Join(args, " ")

	if strings.Contains(cmd, "curl") && strings.Contains(cmd, "|") {
		return Action{true, "curl pipe to shell", args}
	}

	if strings.Contains(cmd, "wget") && strings.Contains(cmd, "&&") {
		return Action{true, "wget chained execution", args}
	}

	if strings.Contains(cmd, ".sh") {
		return Action{true, "shell script execution", args}
	}

	return Action{false, "", args}
}

func RunSandbox(cmd []string) sandbox.Result {
	return sandbox.Run(cmd)
}

func PassThrough(args []string) {
	binary, err := exec.LookPath(args[1])
	if err != nil {
		fmt.Println("Command not found:", args[1])
		return
	}

	c := exec.Command(binary, args[2:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Run()
}
