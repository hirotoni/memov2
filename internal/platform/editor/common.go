package editor

import (
	"fmt"
	"os/exec"
)

type EditorOpener interface {
	Open(basedir, path string) error
}

type DefaultEditorOpener struct{}

var DEO = DefaultEditorOpener{}

func (eo DefaultEditorOpener) Open(basedir, path string) error {
	cmd := exec.Command("code", path, "--folder-uri", basedir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}
	return nil
}
