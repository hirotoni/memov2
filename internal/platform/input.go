package platform

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hirotoni/memov2/internal/common"
)

// ReadLine reads a line interactively from the terminal.
// It opens /dev/tty directly so that it works even when stdin is a pipe.
func ReadLine(prompt string) (string, error) {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return "", common.Wrap(err, common.ErrorTypeService, "failed to open /dev/tty")
	}
	defer tty.Close()

	if prompt != "" {
		fmt.Print(prompt)
	}
	scanner := bufio.NewScanner(tty)
	if !scanner.Scan() {
		return "", common.New(common.ErrorTypeValidation, "canceled")
	}
	if scanner.Err() != nil {
		return "", common.Wrap(scanner.Err(), common.ErrorTypeService, "error reading input")
	}
	return scanner.Text(), nil
}
