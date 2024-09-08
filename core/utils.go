package engine

// Game utils for poker

import (
	"sort"
	"strings"

	msgpb "go-pk-server/gen"
)

func combinations(cards []*msgpb.Card, n int) [][]*msgpb.Card {
	var result [][]*msgpb.Card
	var comb func(start int, chosen []*msgpb.Card)
	comb = func(start int, chosen []*msgpb.Card) {
		if len(chosen) == n {
			combination := make([]*msgpb.Card, n)
			copy(combination, chosen)
			result = append(result, combination)
			return
		}
		for i := start; i <= len(cards)-1; i++ {
			comb(i+1, append(chosen, cards[i]))
		}
	}
	comb(0, []*msgpb.Card{})
	return result
}

func countValues(cards []*msgpb.Card) map[int]int {
	valueCount := make(map[int]int)
	for _, card := range cards {
		valueCount[int(card.Rank)]++
	}
	return valueCount
}

func isFlush(cards []*msgpb.Card) bool {
	firstSuit := cards[0].Suit
	for _, card := range cards {
		if card.Suit != firstSuit {
			return false
		}
	}
	return true
}

func isStraight(cards []*msgpb.Card) bool {
	values := []int{}
	for _, card := range cards {
		values = append(values, int(card.Rank))
	}

	// Sort the values
	sort.Ints(values)

	// Check for consecutive values
	for i := 1; i < len(values); i++ {
		if values[i] != values[i-1]+1 {
			// Special case for A-2-3-4-5 straight (Ace as 1)
			if values[0] == 2 && values[1] == 3 && values[2] == 4 && values[3] == 5 && values[4] == 14 {
				return true
			}
			return false
		}
	}
	return true
}

// getSortedValues returns the sorted values of the cards in descending order
func getSortedValues(cards []*msgpb.Card) []msgpb.RankType {
	// Create a slice to hold the rank values
	ranks := make([]msgpb.RankType, len(cards))

	// Extract the rank from each card
	for i, card := range cards {
		ranks[i] = card.Rank
	}

	// Sort the ranks slice
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i] > ranks[j] // Ascending order
	})

	// Return the sorted ranks
	return ranks
}

// removes repeated elements from the input slice and sorts it in descending order
func removeRepeatedAndSortDesc(cards []*msgpb.Card, elements []msgpb.RankType) []msgpb.RankType {
	// Step 1: Count the number of occurrences of each element
	filter := make(map[msgpb.RankType]bool, len(elements))
	for _, element := range elements {
		filter[element] = true
	}

	// Step 2: Filter out repeated elements
	unique := []msgpb.RankType{}
	for _, value := range cards {
		if !filter[value.Rank] {
			unique = append(unique, value.Rank)
		}
	}

	// Step 3: Sort the result in descending order
	sort.Slice(unique, func(i, j int) bool {
		return unique[i] > unique[j]
	})

	// Return the result
	elements = append(elements, unique...)

	return elements
}

func findHighestCard(cards []*msgpb.Card) msgpb.RankType {
	highest := cards[0].Rank
	for _, card := range cards {
		if card.Rank > highest {
			highest = card.Rank
		}
	}
	return highest
}

func evaluateHand(cards []*msgpb.Card) (msgpb.HankRankingType, []msgpb.RankType) {
	valueCount := countValues(cards)

	// Check for flush and straight
	flush := isFlush(cards)
	straight := isStraight(cards)

	// Royal FLUSH or STRAIGHT FLUSH
	if flush && straight {
		if cards[0].Rank == msgpb.RankType_TEN &&
			cards[1].Rank == msgpb.RankType_TEN &&
			cards[2].Rank == msgpb.RankType_QUEEN &&
			cards[3].Rank == msgpb.RankType_KING &&
			cards[4].Rank == msgpb.RankType_ACE {
			return msgpb.HankRankingType_ROYAL_FLUSH, nil
		}
		return msgpb.HankRankingType_STRAIGH_FLUSH, nil
	}

	// Four of a Kind
	for _, count := range valueCount {
		if count == 4 {
			return msgpb.HankRankingType_FOUR_OF_A_KIND, nil
		}
	}

	// Full House
	three := -1
	pair := -1
	for value, count := range valueCount {
		if count == 3 {
			three = value
		} else if count == 2 {
			pair = value
		}
	}
	if three != -1 && pair != -1 {
		return msgpb.HankRankingType_FULL_HOUSE, nil
	}

	// FLUSH
	if flush {
		return msgpb.HankRankingType_FLUSH, []msgpb.RankType{findHighestCard(cards)}
	}

	// STRAIGHT
	if straight {
		return msgpb.HankRankingType_STRAIGHT, []msgpb.RankType{findHighestCard(cards)}
	}

	// Three of a Kind
	if three != -1 {
		return msgpb.HankRankingType_THREE_OF_A_KIND, removeRepeatedAndSortDesc(cards,
			[]msgpb.RankType{msgpb.RankType(int32(three))})
	}

	// Two Pair
	pairs := []int{}
	for value, count := range valueCount {
		if count == 2 {
			pairs = append(pairs, value)
		}
	}
	if len(pairs) == 2 {
		sort.Sort(sort.Reverse(sort.IntSlice(pairs)))
		return msgpb.HankRankingType_TWO_PAIR, removeRepeatedAndSortDesc(cards,
			[]msgpb.RankType{msgpb.RankType(int32(pairs[0])), msgpb.RankType(int32(pairs[1]))})
	}

	// One Pair
	if len(pairs) == 1 {
		return msgpb.HankRankingType_ONE_PAIR, removeRepeatedAndSortDesc(cards,
			[]msgpb.RankType{msgpb.RankType(int32(pairs[0]))})
	}

	// High Card
	return msgpb.HankRankingType_HIGH_CARD, getSortedValues(cards)
}

func compareTiebreakers(tiebreaker1, tiebreaker2 []msgpb.RankType) int {
	for i := range tiebreaker1 {
		if tiebreaker1[i] < tiebreaker2[i] {
			return 1
		} else if tiebreaker1[i] > tiebreaker2[i] {
			return -1
		}
	}
	return 0
}

func getWinnersName(winners []Player) string {
	// Create a slice to hold the names
	names := make([]string, len(winners))

	// Collect the names
	for i, winner := range winners {
		names[i] = winner.Name()
	}

	// Join the names with ", "
	return strings.Join(names, ", ")
}
