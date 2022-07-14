package main

import (
	"fmt"
	"math/rand"

	tm "github.com/buger/goterm"
	"github.com/gookit/color"
)

type Game struct {
	player_quant int
	//dificulty   int
	table   [63]int
	players [4]Player
}

type Player struct {
	letter   string
	position int
	turn     int
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

func print_table(g *Game) {
	for index := range g.table {
		index += 1
		if index >= 10 {
			color.BgRed.Print(index)
		} else {
			color.BgRed.Print("0", index)
		}
		color.Print(" ")
		if index != 0 && index%10 == 0 {
			fmt.Print("\n\n")
		}
	}
}

func init_players(g *Game) {
	name := ""
	for i := 0; i < g.player_quant; i++ {
		fmt.Printf("\t Player %X letter: ", i+1)
		fmt.Scan(&name)
		g.players[i].letter = name
		g.players[i].turn = i
		g.players[i].position = 0
	}
}

func clear_terminal() {
	tm.Clear()
	tm.MoveCursor(0, 0)
	tm.Flush()
}

func init_game(g *Game) {
	clear_terminal()

	//Start the game wiht an advice =>
	color.BgYellow.Println(" The game start... ")

	//Get configuration (player_quant, dificulty) =>
	get_player_quant(g)
	//TO_DO => get_dificulty(g)

	//init table and players =>
	init_table(g)
	init_players(g)

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
