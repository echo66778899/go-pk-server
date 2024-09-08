package engine

import (
	msgpb "go-pk-server/gen"
	"strings"
)

// Hand represents a hand of playerCards.
type Hand struct {
	playerCards    [2]*msgpb.Card
	rank           msgpb.HankRankingType
	bestHand       []*msgpb.Card
	bestTiebreaker []msgpb.RankType
}

func (h Hand) String() string {
	return h.playerCards[0].String() + ", " + h.playerCards[1].String()
}

func (h *Hand) SetCard(card *msgpb.Card, idx int) {
	h.playerCards[idx] = card
}

func (h *Hand) Cards() []*msgpb.Card {
	return h.playerCards[:]
}

func (h *Hand) Reset() {
	h.playerCards = [2]*msgpb.Card{}
	h.rank = msgpb.HankRankingType(-1)
	h.bestHand = []*msgpb.Card{}
	h.bestTiebreaker = []msgpb.RankType{}
}

func (h *Hand) Evaluate(cc *CommunityCards) (bestRank msgpb.HankRankingType) {
	// Find all possible combinations of 5 playerCards from the hand and community playerCards
	bestRank = msgpb.HankRankingType(-1)
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
	h.rank = bestRank
	return
}

func (h *Hand) BestHand() []*msgpb.Card {
	return h.bestHand
}

func (h *Hand) BestHandString() string {
	var cardsString string
	for _, card := range h.bestHand {
		if cardsString != "" {
			cardsString += ", "
		}
		cardsString += card.String()
	}
	return cardsString
}

func (h *Hand) GetPlayerHandRanking(kicker msgpb.RankType) string {
	var kickerString string
	if kicker > msgpb.RankType_UNSPECIFIED_RANK {
		kickerString = " WITH " + kicker.String() + " KICKER"
	}

	rankString := strings.ReplaceAll(h.rank.String(), "_", " ")

	switch h.rank {
	case msgpb.HankRankingType_HIGH_CARD:
		// Example: "HIGH CARD, ACE"
		return rankString + ", " + getRankValueShortForm(h.bestTiebreaker[0])
	case msgpb.HankRankingType_ONE_PAIR:
		// Example: "ONE PAIR, 3S"
		return rankString + ", " + getRankValueShortForm(h.bestTiebreaker[0]) + "S" + kickerString
	case msgpb.HankRankingType_TWO_PAIR:
		// Example: "TWO PAIR, KINGS AND 3S"
		return rankString + ", " + getRankValueShortForm(h.bestTiebreaker[0]) + "S AND " + h.bestTiebreaker[1].String() + "S" + kickerString
	case msgpb.HankRankingType_THREE_OF_A_KIND:
		// Example: "THREE OF A KIND, 3S" , or "THREE OF A KIND, 3S WITH ACE KICKER"
		return rankString + ", " + getRankValueShortForm(h.bestTiebreaker[0]) + "S" + kickerString
	case msgpb.HankRankingType_STRAIGHT:
		// Example: "TEN HIGH STRAIGHT"
		return getRankValueShortForm(h.bestTiebreaker[0]) + " HIGH " + rankString
	case msgpb.HankRankingType_FLUSH:
		// Example: "KING HIGH FLUSH, SPADES"
		return getRankValueShortForm(h.bestTiebreaker[0]) + " HIGH " + rankString + ", " + h.bestHand[0].Suit.String()
	case msgpb.HankRankingType_FULL_HOUSE:
		// Example: "FULL HOUSE OF KINGS, AND 3S"
		return rankString + " OF " + getRankValueShortForm(h.bestTiebreaker[0]) + "S, AND " + h.bestTiebreaker[1].String() + "S"
	case msgpb.HankRankingType_FOUR_OF_A_KIND:
		// Example: "FOUR OF A KIND OF QUEENS"
		return rankString + " OF " + getRankValueShortForm(h.bestTiebreaker[0]) + "S"
	case msgpb.HankRankingType_STRAIGH_FLUSH:
		// Example: "8 HIGH STRAIGHT FLUSH, DIAMONDS"
		return getRankValueShortForm(h.bestTiebreaker[0]) + " HIGH " + rankString + ", " + h.bestHand[0].Suit.String()
	case msgpb.HankRankingType_ROYAL_FLUSH:
		// Example: "ROYAL FLUSH, SPADES"
		return rankString + ", " + h.bestHand[0].Suit.String()
	default:
		return "Unknown"
	}

}

// Returns the sorted values of the best hand cards in descending order
func (h *Hand) SortedDecendingRankValue() []msgpb.RankType {
	return h.bestTiebreaker
}

func (h *Hand) Compare(otherHand *Hand) int {
	if h.rank > otherHand.rank {
		return 1
	} else if h.rank < otherHand.rank {
		return -1
	}

	return compareTiebreakers(h.bestTiebreaker, otherHand.bestTiebreaker)
}
