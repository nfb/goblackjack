package main

import (
	"math/rand"
	"slices"
	"strconv"
)

// order of suits
var CardSuitId map[string]int = map[string]int{"club": 0, "diamond": 1, "heart": 2, "spade": 3}
var CardSuitAnsi map[string]string = map[string]string{"club": "♣", "diamond": "♦", "heart": "♥", "spade": "♠"}
var pictureCards []string = []string{"A", "J", "Q", "K"}
var suits []string = []string{"club", "diamond", "heart", "spade"}

//var CardRankList []string = ["A","2","3",

type CardDeck struct {
	Cards []ACard
}

type ACard struct {
	Suit        string // club,diamond,heart,spade
	PictureCard bool
	PictureType string // A: ace, J: jack, Q: queen, K: king
	NumericRank int    // Only to be used for numeric cards
}

func (c *ACard) String() string {
	out_suit := CardSuitAnsi[c.Suit]
	if c.PictureCard {
		return c.PictureType + out_suit
	}
	return strconv.Itoa(c.NumericRank) + out_suit
}

// suit - heart, diamond, spade, clubs
// rank - ace, [2-10], jack, queen, king

func genacard() ACard {
	var c ACard
	c.Suit = "spade"
	c.PictureCard = true
	c.PictureType = "A"
	return c
}

func genNumberCards(suit string) []ACard {
	var clist []ACard
	for i := 2; i <= 10; i++ {
		var c ACard
		c.Suit = suit
		c.PictureCard = false
		c.NumericRank = i
		clist = append(clist, c)
	}
	return clist
}

func genPictureCards(suit string) []ACard {
	var clist []ACard
	for _, name := range pictureCards {
		var c ACard
		c.Suit = suit
		c.PictureCard = true
		c.PictureType = name
		clist = append(clist, c)
	}
	return clist
}

func gendeck() CardDeck {
	var cd CardDeck
	for _, s := range suits {
		//cd.Cards = slices.Concat(cd.Cards, genPictureCards(s))
		cd.Cards = append(cd.Cards, genPictureCards(s)...)
		cd.Cards = append(cd.Cards, genNumberCards(s)...)
	}
	return cd
}

func (d *CardDeck) ShuffleCard() {
	deckSize := len(d.Cards)
	opos := rand.Intn(deckSize - 1)
	npos := rand.Intn(deckSize - 1)
	tc := d.Cards[opos]
	interdeck := slices.Delete(d.Cards, opos, opos+1)
	interdeck = slices.Insert(interdeck, npos, tc)
	d.Cards = interdeck
}

func (d *CardDeck) ShuffleLots() {
	for i := 0; i < 208; i++ {
		d.ShuffleCard()
	}
}

func BlackJackHandValues(hand []ACard) []int {
	acecount := 0
	nonacecount := 0
	var handcounts []int

	for _, c := range hand {
		if c.PictureCard {
			if c.PictureType == "A" {
				acecount += 1
			} else {
				nonacecount += 10
			}
		} else {
			nonacecount += c.NumericRank
		}
	}
	for i := 0; i <= acecount; i++ {
		tmpcount := nonacecount
		tmpcount += i * 11
		tmpcount += (acecount - i)
		handcounts = append(handcounts, tmpcount)
	}
	return handcounts
}

func BlackJackHandBestValue(hand []ACard) int {
	var hival int
	busthandbest := 33
	vals := BlackJackHandValues(hand)

	for _, count := range vals {
		if count <= 21 && count > hival {
			hival = count
		}
		if count > 21 && count < busthandbest {
			busthandbest = count
		}
	}
	if hival == 0 {
		return busthandbest
	}
	return hival
}
