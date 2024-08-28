package engine

// Hand represents a hand of cards.
type Hand struct {
	cards          [2]Card
	calRank        HandRank
	bestHand       []Card
	bestTiebreaker []int
}

func (h Hand) String() string {
	return h.cards[0].String() + ", " + h.cards[1].String()
}

func (h *Hand) SetCard(card Card, idx int) {
	h.cards[idx] = card
}

func (h *Hand) Reset() {
	h.cards = [2]Card{}
	h.calRank = HandRank(-1)
	h.bestHand = []Card{}
	h.bestTiebreaker = []int{}
}

func (h *Hand) CalcBestHand(cc *CommunityCards) (bestRank HandRank) {
	// Find all possible combinations of 5 cards from the hand and community cards
	bestRank = HandRank(-1)
	allCards := append(h.cards[:], cc.Cards...)
	allCombinations := combinations(allCards, 5)

	// Find the best hand rank among all combinations
	for _, comb := range allCombinations {
		rank, tiebreaker := evaluateHand(comb)
		if rank > bestRank || (rank == bestRank && compareTiebreakers(tiebreaker, h.bestTiebreaker) > 0) {
			bestRank = rank
			h.bestTiebreaker = tiebreaker
			h.bestHand = comb
		}
	}
	h.calRank = bestRank
	return
}

func (h *Hand) BestHand() string {
	var cardsString string
	for _, card := range h.bestHand {
		cardsString += "(" + card.String() + ")"
	}
	return cardsString
}

func (h *Hand) HandRankString() string {
	return h.calRank.String()
}
