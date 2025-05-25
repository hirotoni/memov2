package platform

import (
	"os"
	"os/exec"
	"strings"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/interfaces"
)

func NewEditor(command string, argsTemplate []string) interfaces.Editor {
	return DefaultEditor{command: command, argsTemplate: argsTemplate}
}

type DefaultEditor struct {
	command      string
	argsTemplate []string
}

func (eo DefaultEditor) Open(basedir, path string) error {
	args := eo.buildArgs(basedir, path)
	cmd := exec.Command(eo.command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error opening editor")
	}
	return nil
}

func (eo DefaultEditor) buildArgs(basedir, path string) []string {
	args := make([]string, len(eo.argsTemplate))
	for i, t := range eo.argsTemplate {
		t = strings.ReplaceAll(t, "{basedir}", basedir)
		t = strings.ReplaceAll(t, "{path}", path)
		args[i] = t
	}
	return args
}
