package termui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type RawTerm struct {
	fd     int
	buf    [1]byte
	state  *term.State
	width  int
	height int
}

func NewRawTerm() *RawTerm {
	return &RawTerm{fd: int(os.Stdin.Fd())}
}

func (rt *RawTerm) Open() error {
	if !term.IsTerminal(rt.fd) {
		return fmt.Errorf("[RawTerm] %s", "stdin if not TTY")
	}
	state, err := term.MakeRaw(rt.fd)
	if err != nil {
		return fmt.Errorf("[RawTerm] %w", err)
	}
	w, h, err := term.GetSize(rt.fd)
	if err != nil {
		return fmt.Errorf("[RawTerm] %w", err)
	}
	rt.state = state
	rt.width = w
	rt.height = h
	return nil
}

func (rt *RawTerm) Close() {
	if rt.state == nil {
		return
	}
	term.Restore(rt.fd, rt.state)
}

func (rt *RawTerm) ReadCh() (byte, error) {
	_, err := os.Stdin.Read(rt.buf[:])
	if err != nil {
		return 0, err
	}
	return rt.buf[0], nil
}

func (rt *RawTerm) ReadLine() (string, error) {
	var out []byte
	for {
		_, err := os.Stdin.Read(rt.buf[:])
		if err != nil {
			return "", err
		}

		ch := rt.buf[0]
		if ch == 13 {
			fmt.Print("\r\n")
			break
		}

		if ch == 127 {
			if len(out) > 0 {
				out = out[:len(out)-1]
				fmt.Print("\b \b")
			}
			continue
		}
		out = append(out, ch)
		fmt.Printf("%c", ch)
	}
	return string(out), nil
}

func (rt *RawTerm) Clear() {
	fmt.Print("\033[H\033[2J")
}

func (rt *RawTerm) Write(format string, a ...any) {
	fmt.Printf(format, a...)
}
