package main

import (
	"bufio"
	"fmt"
	"os"

	libbj "github.com/nfb/goblackjack/pkg"
)

type RoundInstance struct {
	Round *libbj.BlackJackRound
}

func (instance *RoundInstance) GameLoop() {
	reader := bufio.NewReader(os.Stdin)
	var input string
	for {

		instance.Round.StateTransit(input)
		instance.Round.PrintGameState()
		key, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		input = key[:len(key)-1]
	}

}

func main() {
	fmt.Println("vim-go")
	//fmt.Println(CardSuit["diamond"])
	//c := gencard()
	//fmt.Println(c.String())
	d := libbj.Gendeck()
	d.ShuffleLots()
	ri := RoundInstance{}
	ri.Round = &libbj.BlackJackRound{}
	ri.Round.Cards = d.Cards
	ri.GameLoop()
}
