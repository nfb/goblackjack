package pkg

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type BlackJackRound struct {
	Cards            []ACard
	PlayerHand       []ACard
	DealerHand       []ACard
	dealerBottomCard ACard
	PlayerBet        int
	State            int
	nextPlayerQuery  string
	// (pregame)     0:  require bet (transit to either 1 or 2 via valid bet),
	// (first deal)  10: first Cards dealt normal, 11: player blackjack
	// (second deal) 20: post hit not bust, 21: player 21 (non blackjack)
	// (end State)   30: win, 31: push, 32: loose, 33: blackjack
}

// actions values map as follows 0 nilbet?, 1 hit, 2 stand, 3 split, 4 double
// the key represents the game State
var StateActions map[int][]int = map[int][]int{
	0:  {0},
	10: {1, 2},
	11: {1, 2},
	20: {1, 2},
	21: {1, 2},
	30: {},
	31: {},
	32: {},
	33: {},
}

// mutliple of PlayerBet to add to PlayerBet to retuurn to player after round given end State
var StateBetChangeM map[int]float64 = map[int]float64{
	30: 1,
	31: 0,
	32: -1,
	33: 1.5,
}
var ActionName map[int]string = map[int]string{
	0: "bet",
	1: "hit",
	2: "stand",
	3: "split",
	4: "double",
}

func (bj *BlackJackRound) DealPlayer() {
	deckSize := len(bj.Cards)
	bj.PlayerHand = append(bj.PlayerHand, bj.Cards[deckSize-1])
	bj.Cards = slices.Delete(bj.Cards, deckSize-1, deckSize) // should probably use a deckheight variable to reduce the copies going on, fun microbenchmark?
}

func (bj *BlackJackRound) DealDealer() {
	deckSize := len(bj.Cards)
	bj.DealerHand = append(bj.DealerHand, bj.Cards[deckSize-1])
	bj.Cards = slices.Delete(bj.Cards, deckSize-1, deckSize) // should probably use a deckheight variable to reduce the copies going on, fun microbenchmark?
}

func (bj *BlackJackRound) DealHand() {
	bj.DealPlayer()
	bj.DealDealer()
	bj.DealPlayer()
	bj.DealDealer()
	bj.dealerBottomCard = bj.DealerHand[0]
	bj.DealerHand = bj.DealerHand[1:]
}

func (bj *BlackJackRound) PlayDealer() {
	bj.DealerHand = append(bj.DealerHand, bj.dealerBottomCard)
	for BlackJackHandBestValue(bj.DealerHand) < 17 {
		bj.DealDealer()
	}
}

func (bj *BlackJackRound) UpdateState() {
	PlayerCardCount := len(bj.PlayerHand)
	DealerCardCount := len(bj.DealerHand)
	PlayerHandValue := BlackJackHandBestValue(bj.PlayerHand)
	DealerHandValue := BlackJackHandBestValue(bj.DealerHand)
	if PlayerCardCount == 0 {
		bj.State = 0
	}
	if PlayerCardCount == 2 && DealerCardCount == 1 {
		// cards have just been dealt, either blackjack or not
		if PlayerHandValue == 21 {
			bj.State = 11 // Player Blackjack
		} else {
			bj.State = 10 // Non Blackjack first deal
		}
	}
	if PlayerCardCount > 2 && DealerCardCount == 1 {
		// Either busted ourself, got 21 or neither
		if PlayerHandValue > 21 {
			bj.State = 32 // bust
		} else if PlayerHandValue == 21 {
			bj.State = 21
		} else {
			bj.State = 20
		}
	}
	if DealerCardCount > 1 {
		// Player Bust State will be dealt with before this so assumptions made but may need to be tightened up for going into multiple hands

		//standard rules
		if PlayerHandValue < DealerHandValue {
			bj.State = 32
		} else if PlayerHandValue > DealerHandValue {
			bj.State = 30
		} else if PlayerHandValue == DealerHandValue {
			bj.State = 31
		}
		// Busted Exceptions
		if DealerHandValue > 21 {
			bj.State = 30
		} else if PlayerHandValue > 21 { // catch uncaught bust after dealer deals
			bj.State = 32
		}
		// Blackjack exceptions
		if PlayerHandValue == 21 {
			// matching 21 vs 21 not bj vs bj push dealt with in standard rules instance
		} else if DealerHandValue == 21 && DealerCardCount == 2 && PlayerCardCount == 2 { // player blackjack dealer 21
			bj.State = 31
			if DealerHandValue == 21 && DealerCardCount == 2 && PlayerCardCount > 2 { // dealer blackjack player 21
				bj.State = 32
			} else if DealerHandValue == 21 && DealerCardCount > 2 && PlayerCardCount == 2 { // player blackjack dealer 21
				bj.State = 33
			} else if PlayerCardCount == 2 { // blackjack vs normal win
				bj.State = 33
			}
		}
	}
}

func (bj *BlackJackRound) Play(action int) error {
	// -ve input is inverse bet size
	if slices.Contains(StateActions[bj.State], action) {
		if action == 1 { // hit
			bj.DealPlayer()
		} else if action == 2 { // stand
			bj.PlayDealer()
		} else if action < 0 { // bet
		} else {
			return errors.New("Unhandled action")
		}
	} else if action < 0 && slices.Contains(StateActions[bj.State], 0) {
		bj.PlayerBet = -action
		bj.DealHand()
	} else {
		return errors.New("invalid action")
	}
	bj.UpdateState()
	return nil
}

func (bj *BlackJackRound) StateTransit(input string) {
	switch bj.State {
	case 0:
		bj.TakeBet(input)
		if bj.PlayerBet > 0 { // can only advance if valid bet is placed, if any other reason a bet hasn't been set player will retry on next gameloop
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

func (bj *BlackJackRound) SetPostDealState() {
	// add in insurance buy offer on ace card 1
	// add in player blackjack check
	bj.State = 10
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
		bj.PlayerBet += val
	}
	return nil
}

func (bj *BlackJackRound) HandleNormCards(input string) {
	switch input {
	case "i":
		bj.DealPlayer()
		if BlackJackHandBestValue(bj.PlayerHand) > 21 {
			bj.State = 21
			GameEnd("Bust")
		} else {
			bj.State = 20
		}
	case "o":
		bj.DealDealer()
		DealerHand := BlackJackHandBestValue(bj.DealerHand)
		if DealerHand > 21 {
			bj.State = 33
			bj.nextPlayerQuery = "Dealerbust u win"
		}
		res := BlackJackHandBestValue(bj.PlayerHand) - DealerHand
		if res > 0 {
			bj.State = 30
			bj.nextPlayerQuery = "You beat the dealer!"
		} else if res == 0 {
			bj.State = 31
			bj.nextPlayerQuery = "Push, you draw"
		} else {
			bj.State = 32
			bj.nextPlayerQuery = "You loose"
		}
	default:
	}
}
