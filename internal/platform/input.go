package platform

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/hirotoni/memov2/internal/common"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

// ReadLine reads a line interactively from the terminal.
// It opens /dev/tty directly so that it works even when stdin is a pipe.
// Uses raw mode to correctly handle backspace for multibyte characters.
func ReadLine(prompt string) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", common.Wrap(err, common.ErrorTypeService, "failed to open /dev/tty")
	}
	defer tty.Close()

	fd := int(tty.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", common.Wrap(err, common.ErrorTypeService, "failed to set raw mode")
	}
	defer term.Restore(fd, oldState)

	if prompt != "" {
		fmt.Fprint(tty, prompt)
	}

	var buf []rune
	readBuf := make([]byte, 4) // max UTF-8 bytes per rune
	var utf8Buf []byte

	for {
		n, err := tty.Read(readBuf[:1])
		if err != nil || n == 0 {
			fmt.Fprint(tty, "\r\n")
			return "", common.New(common.ErrorTypeValidation, "canceled")
		}

		b := readBuf[0]

		switch {
		case b == 3: // Ctrl+C
			fmt.Fprint(tty, "\r\n")
			return "", common.New(common.ErrorTypeValidation, "canceled")
		case b == 4: // Ctrl+D
			fmt.Fprint(tty, "\r\n")
			return "", common.New(common.ErrorTypeValidation, "canceled")
		case b == 13 || b == 10: // Enter
			fmt.Fprint(tty, "\r\n")
			return string(buf), nil
		case b == 127 || b == 8: // Backspace / Delete
			if len(buf) > 0 {
				removed := buf[len(buf)-1]
				buf = buf[:len(buf)-1]
				w := runewidth.RuneWidth(removed)
				// Move cursor back, overwrite with spaces, move back again
				fmt.Fprint(tty, strings.Repeat("\b", w)+strings.Repeat(" ", w)+strings.Repeat("\b", w))
			}
		case b < 32: // other control characters, ignore
			continue
		default: // printable character (possibly multibyte)
			utf8Buf = append(utf8Buf[:0], b)
			for !utf8.FullRune(utf8Buf) {
				n, err = tty.Read(readBuf[:1])
				if err != nil || n == 0 {
					fmt.Fprint(tty, "\r\n")
					return "", common.New(common.ErrorTypeValidation, "canceled")
				}
				utf8Buf = append(utf8Buf, readBuf[0])
			}
			r, _ := utf8.DecodeRune(utf8Buf)
			buf = append(buf, r)
			fmt.Fprint(tty, string(r))
		}
	}
}
