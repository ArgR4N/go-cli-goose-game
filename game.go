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
func get_player_quant(g *Game) {
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

func init_table(g *Game) {
	for index := range g.table {
		g.table[index] = get_random_case()
	}
}

func init_players(g *Game) {
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

func init_game(g *Game) {
	clear_terminal()

	//Get configuration (player_quant, dificulty) =>
	color.BgHiYellow.Println("Game configuration: ")
	get_player_quant(g)
	//TO_DO => get_dificulty(g)

	//init table and players =>
	init_table(g)
	init_players(g)
}

//Init functions <=

//Game loop functions =>
func print_case(index int, num int) {
	if index+1 >= 10 {
		if num != 0 {
			color.BgRed.Printf(" %d ", index+1)
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

func print_table(g *Game) {
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

func turns(op string, g *Game) {
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
	//color.BgGreen.Print("\nPress enter to throw the dice...")
	//fmt.Scanln()
	time.Sleep(3 * time.Millisecond)
	clear_terminal()
	return n
}

func move_player(dice int, p *Player) {
	p.position += dice
}

func game_loop(g *Game) {
	clear_terminal()
	for g.winner == "" {
		//Turns =>
		turns("sum", g)
		for i := 0; i < g.player_quant; i++ {
			if g.players[i].turn < 1 {
				continue
			}
			//Each active turn =>
			color.BgHiYellow.Printf(" Turn of %s  \n\n", g.players[i].letter)
			print_table(g)       //Print table
			dice := throw_dice() //Throw a dice
			color.BgHiYellow.Printf(" The dice was %o |", dice)

			move_player(dice, &g.players[i]) //Move correct player
			//Aply special cases (turns +/- 1 || moves +3/-2)

			if g.players[i].position > 62 { //Check passed
				move_player(-5, &g.players[i])
			} else if g.players[i].position == 62 { //Check win
				color.Green.Println("\nThe winner is ", g.players[i].letter, "!!")
				g.winner = g.players[i].letter
				break
			}
		}
		turns("sub", g)
	}
}

//Game loop functions <=

//For reset =>
func ask_for_reset(sp *bool) {
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
		game := Game{}                //Create new game
		init_game(&game)              //Init game values
		game_loop(&game)              //Start game loop
		ask_for_reset(&still_playing) //Ask for reset => end program or reset
	}
}
