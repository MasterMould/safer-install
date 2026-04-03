package main

import (
	"fmt"
	"os"
	"strings"

	"safe-install/analyzer"
	"safe-install/interceptor"
	"safe-install/packager"
)

func main() {
	if len(os.Args) < 2 {
	printHelp()
	return
}

// Handle flags
switch os.Args[1] {
case "--help", "-h":
	printHelp()
	return
case "--version":
	fmt.Println("safe-install v1.0.0")
	return
}

// 🌐 Detect raw URL usage (THIS is your missing piece)
if isURL(os.Args[1]) {
	handleURL(os.Args[1:])
	return
}


	func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
	}

	func handleURL(args []string) {
	url := args[0]

	fmt.Println("🌐 Detected URL input")
	fmt.Println("")
	fmt.Println("It looks like you're trying to use:")
	fmt.Println(" ", url)
	fmt.Println("")

	fmt.Println("💡 safe-install works with commands, not raw URLs.")
	fmt.Println("")

	fmt.Println("👉 Try one of these:")
	fmt.Println("")
	fmt.Println("1. Inspect the script safely:")
	fmt.Println("   curl -fsSL", url)
	fmt.Println("")
	fmt.Println("2. Run safely through safe-install:")
	fmt.Println("   safe-install curl -fsSL", url, "| bash")
	fmt.Println("")
	fmt.Println("3. If this is a repository, follow its official install steps.")
	fmt.Println("")
	fmt.Println("❌ No action taken.")
	}

	func printHelp() {
	fmt.Println("🛡️  safe-install — safer script execution")
	fmt.Println("")
	fmt.Println("📦 Usage:")
	fmt.Println("  safe-install <command>")
	fmt.Println("")
	fmt.Println("⚡ Examples:")
	fmt.Println("  safe-install curl -fsSL https://example.com/install.sh | bash")
	fmt.Println("  safe-install wget https://example.com/install.sh && bash install.sh")
	fmt.Println("")
	fmt.Println("🧪 What it does:")
	fmt.Println("  • Detects risky patterns")
	fmt.Println("  • Runs commands in a sandbox")
	fmt.Println("  • Shows system changes")
	fmt.Println("  • Lets you install safely")
	fmt.Println("")
	fmt.Println("🧰 Options:")
	fmt.Println("  --help        Show this help")
	fmt.Println("  --version     Show version")
	fmt.Println("")
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
