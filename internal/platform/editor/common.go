package editor

import (
	"fmt"
	"os/exec"

	"github.com/hirotoni/memov2/internal/config"
	"github.com/hirotoni/memov2/internal/interfaces"
)

func New(c *config.TomlConfig) interfaces.Editor {
	return DefaultEditor{}
}

type DefaultEditor struct{}

func (eo DefaultEditor) Open(basedir, path string) error {
	cmd := exec.Command("code", "--folder-uri", basedir, "--goto", path+":7")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error opening editor: %v", err)
	}
	return nil
}
