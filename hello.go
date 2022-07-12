package main

import "fmt"

type Pair struct{
    a, b interface{}
}

type Player struct{
    name string
    nextTurn int
    position Pair
}


func main(){
  //Here goes the game =>
   fmt.Println("The game start... \n")
   //let a := Game() 
}
