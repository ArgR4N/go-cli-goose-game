package main

import (
	"fmt"
	"math/rand"

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
	table        [63]int
	players      [4]Player
	winner       string
	turn_counter int
}

func clear_terminal() {
	tm.Clear()
	tm.MoveCursor(0, 0)
	tm.Flush()
}

//Init functions =>
func (g *Game) get_player_quant() {
	fmt.Print("Number of player (2 - 4): ")
	n := 0
	fmt.Scan(&n)
	for n < 2 || n > 4 {
		fmt.Print("The number of players must be between 1 and 4: ")
		fmt.Scan(&n)
	}
	g.player_quant = n
}

func get_random_case() int {
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

func (g *Game) init_table() {
	for index := range g.table {
		g.table[index] = get_random_case()
	}
}

func (g *Game) init_players() {
	name := ""
	for i := 0; i < g.player_quant; i++ {
		fmt.Printf("\t Player %X letter: ", i+1)
		fmt.Scan(&name)

		//Init player =>
		g.players[i].letter = name
		g.players[i].turn = 0
		g.players[i].position = 0
	}
}

func (g *Game) init() {
	clear_terminal()

	//Get configuration (player_quant, dificulty) =>
	color.BgHiYellow.Println("Game configuration: ")
	g.get_player_quant()
	//TO_DO => get_dificulty( )

	//init table and players =>
	g.init_table()
	g.init_players()
}

//Init functions <=

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

func (g *Game) print_table() {
	not_player := false
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

func (g *Game) turns(op string) {
	for i := 0; i < g.player_quant; i++ {
		n := -1
		if op == "sum" {
			n = 1
		}
		g.players[i].turn += n
	}
}

func throw_dice() int {
	n := rand.Intn(6) + 1
	color.BgGreen.Print("\nPress enter to throw the dice...")
	fmt.Scanln()
	//time.Sleep(3 * time.Millisecond)
	clear_terminal()
	return n
}

func (g *Game) move_player(dice int, index int) {
	g.players[index].position += dice
}

func (g *Game) apply_special_cases(index int) {
	if g.table[g.players[index].position] != 0 {
		return
	}
}

func (g *Game) start_loop() {
	clear_terminal()
	for g.winner == "" {
		//Turns =>
		g.turns("sum")
		for i := 0; i < g.player_quant; i++ {
			if g.players[i].turn < 1 {
				continue
			}
			color.BgHiYellow.Printf(" Turn of %s  \n\n", g.players[i].letter) //Each active turn =>
			g.print_table()                                                   //Print table
			dice := throw_dice()                                              //Throw a dice
			color.BgHiYellow.Printf(" The dice was %o |", dice)
			g.move_player(dice, i)   //Move correct player
			g.apply_special_cases(i) //Aply special cases (turns +/- 1 || moves +3/-2)

			if g.players[i].position > 62 { //Check passed
				g.move_player(-5, i)
			} else if g.players[i].position == 62 { //Check win
				color.Green.Println("\nThe winner is ", g.players[i].letter, "!!")
				g.winner = g.players[i].letter
				break
			}
		}
		g.turns("sub")
	}
}

//Game loop functions <=

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
