package packager

import (
	"os/exec"
)

func BuildDeb(rootDir string) error {
	cmd := exec.Command(
		"fpm",
		"-s", "dir",
		"-t", "deb",
		"-n", "safe-install-app",
		"-v", "1.0",
		rootDir,
	)

	return cmd.Run()
}
