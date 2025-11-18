package cmd

import (
	"fmt"
	"os"

	"gorogue/internal/termui"
	"gorogue/internal/utils"
)

func Execute() {
	if err := runLoop(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func runLoop() error {
	rawTerm := termui.NewRawTerm()
	err := rawTerm.Open()
	if err != nil {
		return err
	}
	defer rawTerm.Close()

	err = runGame(rawTerm);
	if err != nil {
		return err
	}

	return nil
}

func runGame(rt *termui.RawTerm) error {
	playerName, err := makeInput(rt, "What is your name?")
	if err != nil {
		return nil
	}

	playerClass, _, err := makeMenu(rt, "Choose your class", []string{
		"Paladin",
		"Rogue",
		"Wizard",
	})
	if err != nil {
		return nil
	}

	rt.Write("Welcome %s the %s\r\n", playerName, playerClass)

	return nil
}

const (
	KEY_UP      = 65
	KEY_DOWN    = 66
	KEY_RIGHT   = 67
	KEY_LEFT    = 68
	KEY_CONFIRM = 13
	KEY_BACK    = 127
)

func makeInput(rt *termui.RawTerm, question string) (string, error) {
	rt.Clear()
	rt.Write("[ %s ]\r\n", question)
	rt.Write("\033[31m> \033[0m")
	return rt.ReadLine()
}

func makeMenu(rt *termui.RawTerm, question string, options []string) (string, int, error) {
	cursor := 0

	rt.Clear()
	rt.Write("[ %s ]\r\n", question)
	renderMenu(rt, cursor, options)

	for {
		ch, err := rt.ReadCh()
		if err != nil {
			return "", -1, err
		}
		switch ch {
		case KEY_BACK:
			return "", -1, nil
		case KEY_CONFIRM:
			return options[cursor], cursor, nil
		case KEY_UP:
			cursor = cursor - 1 + len(options)
		case KEY_DOWN:
			cursor = cursor + 1
		}
		cursor %= len(options)
		rt.Clear()
		rt.Write("[ %s ]\r\n", question)
		renderMenu(rt, cursor, options)
	}
}

func renderMenu(rt *termui.RawTerm, cursor int, options []string) {
	for i, opt := range options {
		char := utils.IfElse(i == cursor, ">", " ")
		rt.Write("\033[31m%s\033[0m %s\r\n", char, opt)
	}
}
