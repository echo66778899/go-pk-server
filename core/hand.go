package engine

// Hand represents a hand of playerCards.
type Hand struct {
	playerCards    [2]Card
	itsRank        HandRanking
	bestHand       []Card
	bestTiebreaker []int
}

func (h Hand) String() string {
	return h.playerCards[0].String() + ", " + h.playerCards[1].String()
}

func (h *Hand) SetCard(card Card, idx int) {
	h.playerCards[idx] = card
}

func (h *Hand) Reset() {
	h.playerCards = [2]Card{}
	h.itsRank = HandRanking(-1)
	h.bestHand = []Card{}
	h.bestTiebreaker = []int{}
}

func (h *Hand) Evaluate(cc *CommunityCards) (bestRank HandRanking) {
	// Find all possible combinations of 5 playerCards from the hand and community playerCards
	bestRank = HandRanking(-1)
	allCards := append(h.playerCards[:], cc.Cards...)
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
	h.itsRank = bestRank
	return
}

func (h *Hand) BestHand() string {
	var cardsString string
	for _, card := range h.bestHand {
		if cardsString != "" {
			cardsString += ", "
		}
		cardsString += card.String()
	}
	return cardsString
}

func (h *Hand) HandRankingString() string {
	return h.itsRank.String()
}

func (h *Hand) Kicker() []int {
	return h.bestTiebreaker
}

func (h *Hand) Compare(otherHand *Hand) int {
	if h.itsRank > otherHand.itsRank {
		return 1
	} else if h.itsRank < otherHand.itsRank {
		return -1
	}

	return compareTiebreakers(h.bestTiebreaker, otherHand.bestTiebreaker)
}
