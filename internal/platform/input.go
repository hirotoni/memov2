package platform

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hirotoni/memov2/internal/common"
)

// ReadLine reads a line from standard input
func ReadLine(prompt string) (string, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", common.New(common.ErrorTypeValidation, "canceled")
	}
	if scanner.Err() != nil {
		return "", common.Wrap(scanner.Err(), common.ErrorTypeService, "error reading input")
	}
	return scanner.Text(), nil
}
