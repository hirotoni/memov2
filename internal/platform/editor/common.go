package editor

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hirotoni/memov2/internal/config"
)

type Editor interface {
	Open(basedir, path string) error
}

func New(c *config.TomlConfig) Editor {
	return DefaultEditor{}
}

type DefaultEditor struct{}

var DEO = DefaultEditor{}

func (eo DefaultEditor) Open(basedir, path string) error {
	cmd := exec.Command("code", "--folder-uri", basedir, "--goto", path+":7")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}
	return nil
}

type MockEditor struct{}

func (eo MockEditor) Open(basedir, path string) error {
	_, err := os.Stat(basedir)
	if err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}
	_, err = os.Stat(path)
	if err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}
	return nil
}
