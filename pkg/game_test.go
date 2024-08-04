package pkg

import "testing"

func TestDealTo17(t *testing.T) {
	// Build test deck of 1's
	// Confirm PlayDealer will hit up to 17
	var testRound BlackJackRound
	testRound.DealerHand = append(testRound.DealerHand, ACard{PictureCard: false, NumericRank: 1})
	testRound.dealerBottomCard = ACard{PictureCard: false, NumericRank: 1}
	testRound.PlayerHand = append(testRound.PlayerHand, ACard{PictureCard: false, NumericRank: 1})
	testRound.PlayerHand = append(testRound.PlayerHand, ACard{PictureCard: false, NumericRank: 1})
	for i := 1; i < 24; i++ {
		testRound.Cards = append(testRound.Cards, ACard{PictureCard: false, NumericRank: 1})
	}
	testRound.PlayDealer()
	//testRound.DealHand()
	//testRound.Play(-1)
	//testRound.Play(
	if BlackJackHandBestValue(testRound.DealerHand) < 17 {
		t.Fatalf("BlackJackHandBestValue(testRound.DealerHand) < 17")
	}

}

func TestBlackJackPush(t *testing.T) {
	// test for push scenario if both dealer and player have blackjack
}
