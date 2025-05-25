package components

import (
	"log"
	"os/exec"
)

// OpenEditor opens editor
func OpenEditor(basedir, path string) {
	cmd := exec.Command("code", path, "--folder-uri", basedir)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
