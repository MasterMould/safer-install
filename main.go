package main

import (
	"fmt"
	"os"

	"safe-install/interceptor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: safe-install <command>")
		return
	}

	action := interceptor.Check(os.Args)

	if action.Intercept {
		fmt.Println("⚠️ Risky install detected:", action.Reason)
		fmt.Println("Running in sandbox...\n")

		result := interceptor.RunSandbox(action.Command)

		fmt.Println("Analysis complete:")
		fmt.Println(result.Summary)
	} else {
		interceptor.PassThrough(os.Args)
	}
}
