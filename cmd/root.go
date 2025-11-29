package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"gorogue/internal/termui"
	"gorogue/internal/utils"
)

func Execute() {
	rt := termui.NewRawTerm()
	err := rt.Open()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	defer rt.Close()

	err = runGame(rt)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

type Player struct {
	Name    string
	Class   string
	Hp      int
	Mp      int
	Ap      int
	Potions int
	Skills  []string
}

type Enemy struct {
	Name string
	Hp   int
	Ap   int
}

func runGame(rt *termui.RawTerm) error {
	var err error
	var response string

	player := Player{
		Potions: 10,
		Skills:  []string{"Attack", "Magic", "Potion"},
	}

	response, err = makeInput(rt, "What is your name?")
	if err != nil {
		return nil
	}
	player.Name = response

	response, _, err = makeMenu(rt, "Choose your class", []string{
		"Paladin",
		"Rogue",
		"Wizard",
	})
	if err != nil {
		return err
	}
	player.Class = response

	player.Hp = utils.IfThenElse(player.Class == "Paladin", 120, 80)
	player.Mp = utils.IfThenElse(player.Class == "Wizard", 40, 20)
	player.Ap = utils.IfThenElse(player.Class == "Rogue", 8, 4)

	rt.Write("Welcome %s the %s\r\n", player.Name, player.Class)
	makeNext(rt)
	rt.Clear()

	enemyTypes := []Enemy{}
	enemyTypes = append(enemyTypes, Enemy{
		Name: "Goblin",
		Hp:   10,
		Ap:   12,
	})
	enemyTypes = append(enemyTypes, Enemy{
		Name: "Skeleton",
		Hp:   40,
		Ap:   5,
	})
	enemyTypes = append(enemyTypes, Enemy{
		Name: "Bandit",
		Hp:   30,
		Ap:   4,
	})
	enemiesToKill := 3
	spellDamage := 10
	potionStrenght := 25

	for {
		enemyType := enemyTypes[rand.Intn(len(enemyTypes))]
		enemyCurrentHp := enemyType.Hp
		rt.Write("Enemy %s is approaching...\r\n", enemyType.Name)

		for {
			rt.Write(
				"[%s]\r\nHP: %s | MP: %s | Potions: %s\r\n",
				col(termui.CLR_YELLOW, player.Name),
				col(termui.CLR_RED, strconv.Itoa(player.Hp)),
				col(termui.CLR_BLUE, strconv.Itoa(player.Mp)),
				col(termui.CLR_CYAN, strconv.Itoa(player.Potions)),
			)
			rt.Write(
				"[%s]\r\nHP: %s | AP: %s\r\n",
				col(termui.CLR_YELLOW, enemyType.Name),
				col(termui.CLR_RED, strconv.Itoa(enemyCurrentHp)),
				col(termui.CLR_BRED, strconv.Itoa(enemyType.Ap)),
			)

			makeNext(rt)

			rt.Clear()
			response, _, err = makeMenu(rt, "Your turn", player.Skills)
			if err != nil {
				return err
			}

			rt.Clear()

			switch response {
			case "Attack":
				{
					rt.Write("%s attacks %s for %d points\r\n", col(termui.CLR_YELLOW, player.Name), col(termui.CLR_RED, enemyType.Name), player.Ap)
					enemyCurrentHp -= player.Ap
					break
				}
			case "Magic":
				{
					if player.Mp <= 0 {
						rt.Write("%s fails to cast spell\r\n", col(termui.CLR_YELLOW, player.Name))
					} else {
						rt.Write("%s casts the spell on %s for %d points\r\n", col(termui.CLR_YELLOW, player.Name), col(termui.CLR_RED, enemyType.Name), spellDamage)
						enemyCurrentHp -= spellDamage
						player.Mp -= 5
					}
					break
				}
			case "Potion":
				{
					if player.Potions <= 0 {
						rt.Write("%s has no potions left\r\n", col(termui.CLR_YELLOW, player.Name))
					} else {
						rt.Write("%s drinks potion and heals for %d points\r\n", col(termui.CLR_YELLOW, player.Name), potionStrenght)
						player.Hp += potionStrenght
						player.Potions--
					}
				}
			}

			if enemyCurrentHp <= 0 {
				rt.Write(
					"%s the %s defeated %s\r\n",
					col(termui.CLR_YELLOW, player.Name),
					col(termui.CLR_CYAN, player.Class),
					col(termui.CLR_RED, enemyType.Name),
				)
				makeSeparator(rt)
				makeNext(rt)
				break
			}

			rt.Write("%s attacks %s for %d points\r\n", col(termui.CLR_RED, enemyType.Name), col(termui.CLR_YELLOW, player.Name), enemyType.Ap)
			player.Hp -= enemyType.Ap
		}

		enemiesToKill--
		if enemiesToKill < 1 {
			break
		}
		makeSeparator(rt)
		makeNext(rt)
	}

	return nil
}

func makeSeparator(rt *termui.RawTerm) {
	rt.Write("\r\n----------------------------------------\r\n")
}

func makeNext(rt *termui.RawTerm) error {
	rt.Write("Press ENTER to continue\r\n")
	for {
		ch, err := rt.ReadCh()
		if err != nil {
			return err
		}
		if ch == termui.KEY_CTRL_C {
			return fmt.Errorf("Exit")
		}
		if ch == termui.KEY_CR {
			break
		}
	}
	return nil
}

func makeInput(rt *termui.RawTerm, question string) (string, error) {
	rt.Clear()
	rt.Write("[ %s ]\r\n", col(termui.CLR_YELLOW, question))
	rt.Write("%s", col(termui.CLR_RED, "> "))
	return rt.ReadLine()
}

func makeMenu(rt *termui.RawTerm, question string, options []string) (string, int, error) {
	cursor := 0

	rt.Clear()
	rt.Write("[ %s ]\r\n", col(termui.CLR_YELLOW, question))
	renderMenu(rt, cursor, options)

	for {
		ch, err := rt.ReadCh()
		if err != nil {
			return "", -1, err
		}
		switch ch {
		case termui.KEY_CTRL_C:
			return "", -1, fmt.Errorf("Exit")
		case termui.KEY_CR:
			return options[cursor], cursor, nil
		case termui.KEY_UP:
			cursor = cursor - 1 + len(options)
		case termui.KEY_DOWN:
			cursor = cursor + 1
		}
		cursor %= len(options)
		rt.Clear()
		rt.Write("[ %s ]\r\n", col(termui.CLR_YELLOW, question))
		renderMenu(rt, cursor, options)
	}
}

func renderMenu(rt *termui.RawTerm, cursor int, options []string) {
	for i, opt := range options {
		char := utils.IfThenElse(i == cursor, ">", " ")
		color := utils.IfThenElse(i == cursor, termui.CLR_YELLOW, termui.CLR_RESET)
		rt.Write("%s %s\r\n", col(termui.CLR_RED, char), col(color, opt))
	}
}

func col(color int, message string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, message)
}
