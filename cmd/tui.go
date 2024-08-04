package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	libbj "github.com/nfb/goblackjack/pkg"
)

type RoundInstance struct {
	Round *libbj.BlackJackRound
}

func (instance *RoundInstance) PrintGameState() {
	fmt.Println("---")
	fmt.Print("Dealer Hand: ")
	for _, c := range instance.Round.DealerHand {
		fmt.Print(c.String())
		fmt.Print(" ")
	}
	fmt.Println("")
	fmt.Print("Player Hand: ")
	for _, c := range instance.Round.PlayerHand {
		fmt.Print(c.String())
		fmt.Print(" ")
	}
	fmt.Println("Player Bet: ", strconv.Itoa(instance.Round.PlayerBet))
	fmt.Println("---")
	if instance.Round.State < 30 {
		fmt.Println("Available actions:")
		for _, action := range libbj.StateActions[instance.Round.State] {
			fmt.Println(libbj.ActionName[action])
		}
	} else {
		fmt.Println("Round finished")
		balanceReturned := float64(instance.Round.PlayerBet) + float64(instance.Round.PlayerBet)*libbj.StateBetChangeM[instance.Round.State]
		fmt.Println("Balance returned: ", balanceReturned)
		os.Exit(0)
	}

	//fmt.Println(instance.Round.nextPlayerQuery)
}

func (instance *RoundInstance) GameLoop() {
	reader := bufio.NewReader(os.Stdin)
	var input string
	for {

		instance.PrintGameState()
		key, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		input = key[:len(key)-1]
		switch input {
		case "hit":
			instance.Round.Play(1)
		case "stand":
			instance.Round.Play(2)
		default:
			i, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("I DONT KNOW HOW TO DEAL WITH THIS ACTION")
			} else if i < 1 {
				fmt.Println("You can't try and place a negative bet...")
			} else {
				instance.Round.Play(-i)
			}
		}
	}
}

func main() {
	//fmt.Println(CardSuit["diamond"])
	//c := gencard()
	//fmt.Println(c.String())

	d := libbj.Gendeck()
	d.ShuffleLots()

	ri := RoundInstance{}
	ri.Round = &libbj.BlackJackRound{}
	ri.Round.Cards = d.Cards
	//ri.Round.Cards = libbj.RiggedDeck
	ri.GameLoop()
}
