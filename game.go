package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	tm "github.com/buger/goterm"
	"github.com/gookit/color"
)

var ops = map[string][]string{"Auto dice": {"Y", "N", "I"}, "Player quant": {"2", "3", "4"}}

type Player struct {
	letter   string
	position int
	turn     int
}

type Game struct {
	player_quant int
	//dificulty   int
	table     [63]int
	players   [4]Player
	winner    string
	auto_dice bool
	turn_move bool
}

//Auxiliar functions =>
func clear_terminal() {
	tm.Clear()
	tm.MoveCursor(0, 0)
	tm.Flush()
}

func abs(n int) int {
	if n < 0 {
		return n * -1
	}
	return n
}

func get_sprial_array(arr [63]int) [63]int {
	return [63]int{}
}

//Init functions =>

func random_case() int {
	newCase := rand.Float64()
	if newCase < 0.7 {
		return 0
	} else {
		n := 0
		for n == 0 {
			n = rand.Intn(5) - 2
		}
		return n
	}
}

func (g *Game) init_table() {
	for index := range g.table {
		g.table[index] = random_case()
	}
}

func (g *Game) init_players() {
	for i := 0; i < g.player_quant; i++ {
		name := ""
		fmt.Printf("\t Player %X letter: ", i+1)
		fmt.Scan(&name)
		g.players[i] = Player{
			letter:   name,
			turn:     1,
			position: 0,
		}
	}
}

func (g *Game) get_config(s string) {
	config_print := func(s string) {
		fmt.Print(s, "? =>")
		options := append(ops[s], "I")
		fmt.Print(" [ ") //general beginning
		for index, char := range options {
			fmt.Print(char)
			if index == len(options)-1 {
				continue
			}
			fmt.Print(" / ")
		}
		fmt.Print(" ] : ") //general end
	}
	config_prompt := func(s string) string {
		input := ""
		fmt.Scan(&input)
		options := append(ops[s], "I")
		for i, elm := range options {
			if input == elm {
				break
			} else if i == len(options)-1 {
				input = ""
			}
		}

		return input
	}
	config_show_info := func(input *string) {
		if *input == "I" {
			fmt.Print("Info of the option goes here!")
			*input = ""
		}
	}
	config_set_value := func(s string, input string) {
		//Provisional =>
		if s == "Auto dice" {
			g.auto_dice = input == "Y"
		}
		if s == "Player quant" {
			inp, _ := strconv.Atoi(input)
			g.player_quant = inp
		}
	}
	input := ""
	for input == "" {
		config_print(s)
		input = config_prompt(s)
		config_show_info(&input)
	}
	config_set_value(s, input)
}

func (g *Game) init() {
	clear_terminal()
	color.BgYellow.Println("Game configuration: ")
	g.get_config("Auto dice")
	g.get_config("Player quant")
	//TO DO => g.get_config("Show table")
	//TODO => g.get_config("Difficulty")
	g.init_table()
	g.init_players()
	clear_terminal()
}

//Game loop functions =>
func print_case(index int, num int) {
	output := ""
	if index+1 < 10 {
		output = "0"
	}
	switch num {
	case 1:
		color.BgHiGreen.Printf(" %s%d ", output, index+1)
	case 2:
		color.BgGreen.Printf(" %s%d ", output, index+1)
	case -1:
		color.BgHiRed.Printf(" %s%d ", output, index+1)
	case -2:
		color.BgRed.Printf(" %s%d ", output, index+1)
	default:
		color.BgGray.Printf(" %s%d ", output, index+1)
	}
}

func (g *Game) print_table(dice int, p Player) {
	not_player := false
	if g.turn_move { //Just in 'regular' turns
		color.BgYellow.Printf("The dice is %o | Is the turn of %s \n\n", dice, p.letter)
	} else {
		switch dice < 0 {
		case true:
			color.BgRed.Printf("You move %o cases | Is the turn of %s \n\n", dice, p.letter)
		case false:
			color.BgGreen.Printf("You move %o cases | Is the turn of %s \n\n", dice, p.letter)

		}
	}
	for index, elm := range g.table {
		for i := 0; i < g.player_quant; i++ {
			if g.players[i].position == index && !not_player {
				color.BgYellow.Printf(" %s ", g.players[i].letter)
				color.Print(" ")
				not_player = true
			}
		}
		if !not_player {
			print_case(index, elm)
			color.Print(" ")
		}
		if (index+1)%10 == 0 { //New line
			fmt.Print("\n\n")
		}
		not_player = false
	}
	fmt.Println()
}

func (p *Player) turns(op string) {
	n := 1
	if op == "sub" {
		n = -1
	}
	p.turn += n
}

func throw_dice(auto bool) int {
	n := rand.Intn(6) + 1
	if auto {
		time.Sleep(1 * time.Second)

	} else {
		color.BgGreen.Print("\nPress enter to throw the dice...")
		fmt.Scanln()
	}
	return n
}

func (g *Game) move_player(dice int, index int) {
	if g.players[index].position+dice > 63 {
		color.BgRed.Println("You passed the end for %o cases! ", g.players[index].position+dice-63)
		dice = 63 - g.players[index].position + dice - 5
	}
	for i := 0; i < abs(dice); i++ {
		if i != 0 { //Don´t wait in the first 'step'
			time.Sleep(500 * time.Millisecond)
		}
		clear_terminal()
		g.players[index].position += dice / abs(dice)
		g.print_table(dice, g.players[index])
	}
}

func (g *Game) apply_special_cases(index int) {
	g.turn_move = false
	for !g.turn_move {
		switch g.table[g.players[index].position] {
		case -1:
			g.players[index].turns("sub")
			break
		case -2:
			g.move_player(-2, index)
		case 1:
			dice := throw_dice(true)
			g.move_player(dice, index)
		case 2:
			g.move_player(2, index)
		case g.table[g.players[index].position]:
			break
		}
		g.turn_move = true
	}
}

func (g *Game) start_loop() {
	clear_terminal()
	for g.winner == "" {
		//Turns =>
		for i := 0; i < g.player_quant; i++ {
			g.turn_move = true         // Know if the move is from an normal turn or ar special case
			if g.players[i].turn < 1 { //Check turn
				g.players[i].turns("sum")
				continue
			}
			dice := throw_dice(g.auto_dice)  //Throw a dice
			g.move_player(dice, i)           //Move correct player
			g.apply_special_cases(i)         //Apply special cases
			if g.players[i].position == 62 { //Check win
				//win()
				color.Green.Println("\nThe winner is ", g.players[i].letter, "!! \n")
				g.winner = g.players[i].letter
				break
			}
			if g.players[i].position > 62 { //Check passed
				g.move_player(-6, i)
			}
		}
	}
}

//For reset =>
func (g *Game) ask_for_reset(sp *bool) {
	input := ""
	for input != "Y" && input != "N" {
		color.BgRed.Print("\n Do you want play again?[Y/N]: ")
		fmt.Scan(&input)
	}
	*sp = input == "Y"
}

func main() {
	still_playing := true
	for still_playing {
		game := Game{}                     //Create new game
		game.init()                        //Init game values
		game.start_loop()                  //Start game loop
		game.ask_for_reset(&still_playing) //Ask for reset => End program or Reset
	}
}
