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
	table   [63]int
	players [4]Player
	winner  string
}

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

func move_player(dice int, p *Player) {
	p.position += dice
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

func game_loop(g *Game) {
	for g.winner == "" {
		//Turns =>
		turns("sum", g)
		for i := 0; i < g.player_quant; i++ {
			if g.players[i].turn < 1 {
				continue
			}

			//Each turn=>
			//print_table(g) //Print table
			for i := 0; i < g.player_quant; i++ { //Show info
				fmt.Println(g.players)
			}

			dice := throw_dice()             //Throw a dice
			move_player(dice, &g.players[i]) //Move correct player
			//Aply special cases (turns +/-1 || moves +3/-2)

			if g.players[i].position > 62 { //Check passed
				move_player(-5, &g.players[i])
			} else if g.players[i].position == 62 { //Check win
				g.winner = g.players[i].letter
			}
		}
		turns("sub", g)
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

	//Start game loop =>
	game_loop(g)
}

func ask_for_reset() bool {
	input := ""
	for input != "Y" && input != "N" {
		fmt.Print("\nDo you want play again?[Y/N]: ")
		fmt.Scan(&input)
	}

	return input == "Y"
}

func main() {
	for {
		game := Game{}
		init_game(&game)
		if !ask_for_reset() {
			break
		}
	}
}
