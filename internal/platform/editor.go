package platform

import (
	"os/exec"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/hirotoni/memov2/internal/interfaces"
)

func NewEditor() interfaces.Editor {
	return DefaultEditor{}
}

type DefaultEditor struct{}

func (eo DefaultEditor) Open(basedir, path string) error {
	cmd := exec.Command("code", "--folder-uri", basedir, "--goto", path+":7")
	if err := cmd.Run(); err != nil {
		return common.Wrap(err, common.ErrorTypeService, "error opening editor")
	}
	return nil
}
