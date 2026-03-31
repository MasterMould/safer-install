package main

import (
	"fmt"
	"os"

	"safe-install/analyzer"
	"safe-install/interceptor"
	"safe-install/packager"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: safe-install <command>")
		return
	}

	action := interceptor.Check(os.Args)

	if !action.Intercept {
		interceptor.PassThrough(os.Args)
		return
	}

	fmt.Println("⚠️ Intercepted:", action.Reason)
	fmt.Println("🔬 Running in sandbox...\n")

	result := interceptor.RunSandbox(action.Command)

	score, level := analyzer.Score(result.Diff)

	fmt.Println("===== ANALYSIS =====")
	fmt.Println(result.Diff)
	fmt.Println("--------------------")
	fmt.Printf("Risk Score: %d (%s)\n", score, level)

	fmt.Println("\nOptions:")
	fmt.Println("[1] Build .deb package")
	fmt.Println("[2] Run anyway")
	fmt.Println("[3] Abort")

	var choice string
	fmt.Print("> ")
	fmt.Scanln(&choice)

	switch choice {
	case "1":
		err := packager.BuildDeb(result.RootDir)
		if err != nil {
			fmt.Println("❌ Packaging failed:", err)
		} else {
			fmt.Println("✅ Package created")
		}
	case "2":
		interceptor.PassThrough(os.Args)
	default:
		fmt.Println("❌ Aborted")
	}
}
