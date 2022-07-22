package main

import (
	"fmt"
	"math/rand"
	"time"

	tm "github.com/buger/goterm"
	"github.com/gookit/color"
)

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

//Init functions =>
func (g *Game) get_player_quant() {
	fmt.Print("Number of player (2 - 4): ")
	n := 0
	fmt.Scan(&n)
	for n < 2 || n > 4 {
		fmt.Print("The number of players must be between 1 and 4: ")
		fmt.Scan(&n)
		fmt.Println() //New line
	}
	g.player_quant = n
}

func get_random_case() int {

	return 0

	newCase := rand.Float64()
	if newCase < 0.7 {
		return 0
	} else {
		n := rand.Intn(4) - 2
		for n == 0 {
			n = rand.Intn(4) - 2
		}
		return n
	}
}

func (g *Game) get_auto_dice() {
	input := ""
	for input != "Y" && input != "N" {
		fmt.Print("\nAuto dice? [Y/N]: ")
		fmt.Scan(&input)
	}
	g.auto_dice = input == "Y"
}

func (g *Game) init_table() {
	for index := range g.table {
		g.table[index] = get_random_case()
	}
}

func (p *Player) init_player(name string) {
	p.letter = name
	p.turn = 1
	p.position = 0
}

func (g *Game) init_players() {
	name := ""
	for i := 0; i < g.player_quant; i++ {
		fmt.Printf("\t Player %X letter: ", i+1)
		fmt.Scan(&name)
		g.players[i].init_player(name)
	}
}

func (g *Game) init() {
	clear_terminal()
	color.BgYellow.Println("Game configuration: ")

	//Get configuration (player_quant, dificulty, auto_dice and show_table) =>
	g.get_auto_dice()
	//TO DO => g.get_show_table()
	g.get_player_quant()
	//TODO => get_dificulty( )

	//init table and players =>
	g.init_table()
	g.init_players()
}

//Game loop functions =>
func print_case(index int, num int) {
	if index+1 >= 10 {
		if num != 0 {
			if num < 0 {
				color.BgRed.Printf(" %d ", index+1)
			} else {
				color.BgGreen.Printf(" %d ", index+1)
			}
		} else {
			color.BgGray.Printf(" %d ", index+1)
		}
	} else {
		if num != 0 {
			color.BgRed.Printf(" 0%d ", index+1)
		} else {
			color.BgGray.Printf(" 0%d ", index+1)
		}
	}
}

func (g *Game) print_table(dice int, p Player) {
	not_player := false
	if dice != 0 {
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
		n = 1
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
	//g.apply_special_cases(index)
	//TO DO
	// -1 =>Pierde un turno
	// -2 => Retrocede 2 casillas
	// 1 => Gana un turno
	// 2 => Avanza 3 casillas
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
			dice := throw_dice(g.auto_dice) //Throw a dice
			g.move_player(dice, i)          //Move correct player
			g.apply_special_cases(i)        //Apply special cases
			if g.players[i].position == 62 { //Check win
				//win()
				color.Green.Println("\nThe winner is ", g.players[i].letter, "!! \n")
				g.winner = g.players[i].letter
				g.print_table(0, Player{})
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
		game.ask_for_reset(&still_playing) //Ask for reset => end program or reset
	}
}
