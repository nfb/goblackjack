package pkg

import (
	"fmt"
	"os"
	"slices"
	"strconv"
)

type BlackJackRound struct {
	Cards           []ACard
	playerhand      []ACard
	dealerhand      []ACard
	playerbet       int
	state           int
	nextPlayerQuery string
	// 0: require bet (transit to either 1 or 2 via valid bet),
	// 10: first Cards dealt normal, 11: player blackjack
	// 20: post hit not bust, 21: bust
	// 30: win, 31: push, 32: loose, 33: dealerbust win
}

type BlackJackAction struct {
	Action int // 0 bet, 1 hit, 2 stand, 3 split, 4 double
	Amount int
}

func (bj *BlackJackRound) DealPlayer() {
	deckSize := len(bj.Cards)
	bj.playerhand = slices.Concat(bj.playerhand, bj.Cards[deckSize-1:])
	bj.Cards = bj.Cards[:deckSize-1] // should probably use a deckheight variable to reduce the copies going on, fun microbenchmark?
}

func (bj *BlackJackRound) DealDealer() {
	deckSize := len(bj.Cards)
	bj.dealerhand = slices.Concat(bj.dealerhand, bj.Cards[deckSize-1:])
	bj.Cards = bj.Cards[:deckSize-1] // should probably use a deckheight variable to reduce the copies going on, fun microbenchmark?
}

//func (bj *BlackJackRound) PostDealCheck() {
//	var dealerhandvalues
//	playerhandvalue := BlackJackHandValues(bj.playerhand)
//	if len(dealerhand) {
//

func (bj *BlackJackRound) PrintGameState() {
	fmt.Println("---")
	fmt.Print("Dealer Hand: ")
	for _, c := range bj.dealerhand {
		fmt.Print(c.String())
		fmt.Print(" ")
	}
	fmt.Println("")
	fmt.Print("Player Hand: ")
	for _, c := range bj.playerhand {
		fmt.Print(c.String())
		fmt.Print(" ")
	}
	fmt.Println("Player Bet: ", strconv.Itoa(bj.playerbet))
	fmt.Println("---")
	fmt.Println(bj.nextPlayerQuery)
}

func (bj *BlackJackRound) DealHand() {
	bj.DealPlayer()
	bj.DealDealer()
	bj.DealPlayer()
}

func (bj *BlackJackRound) SetPostDealState() {
	// add in insurance buy offer on ace card 1
	// add in player blackjack check
	bj.state = 10
	bj.nextPlayerQuery = "H(i)t, H(o)ld or (D)ouble"
}

func (bj *BlackJackRound) TakeBet(input string) error {
	if input == "" {
		return nil
	}
	val, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if val > 0 {
		bj.playerbet += val
	}
	return nil
}

func (bj *BlackJackRound) HandleNormCards(input string) {
	switch input {
	case "i":
		bj.DealPlayer()
		if BlackJackHandBestValue(bj.playerhand) > 21 {
			bj.state = 21
			GameEnd("Bust")
		} else {
			bj.state = 20
		}
	case "o":
		bj.DealDealer()
		dealerhand := BlackJackHandBestValue(bj.dealerhand)
		if dealerhand > 21 {
			bj.state = 33
			bj.nextPlayerQuery = "Dealerbust u win"
		}
		res := BlackJackHandBestValue(bj.playerhand) - dealerhand
		if res > 0 {
			bj.state = 30
			bj.nextPlayerQuery = "You beat the dealer!"
		} else if res == 0 {
			bj.state = 31
			bj.nextPlayerQuery = "Push, you draw"
		} else {
			bj.state = 32
			bj.nextPlayerQuery = "You loose"
		}
	default:
	}
}

func (bj *BlackJackRound) StateTransit(input string) {
	switch bj.state {
	case 0:
		bj.TakeBet(input)
		if bj.playerbet > 0 { // can only advance if valid bet is placed, if any other reason a bet hasn't been set player will retry on next gameloop
			bj.DealHand()
			bj.SetPostDealState()
		} else {
			bj.nextPlayerQuery = "enter bet as integer number"
		}
	case 10:
		bj.HandleNormCards(input)
	case 20:
		bj.HandleNormCards(input)

	case 30:
		os.Exit(0)
	case 31:
		os.Exit(0)
	case 32:
		os.Exit(0)
	case 33:
		os.Exit(0)

	default:
		fmt.Println("ummM")
		// gamestart
		// request initial bet
	}
}

func GameEnd(message string) {
	fmt.Println(message)
	os.Exit(0)
}
