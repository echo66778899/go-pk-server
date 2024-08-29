package engine

// Game utils for poker

import "sort"

// Hand ranking constants.
type HandRanking int

const (
	HighCard HandRanking = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (phr HandRanking) String() string {
	return [...]string{"HighCard", "OnePair", "TwoPair", "ThreeOfAKind", "Straight",
		"Flush", "FullHouse", "FourOfAKind", "StraightFlush", "RoyalFlush"}[phr]
}

func combinations(cards []Card, n int) [][]Card {
	var result [][]Card
	var comb func(start int, chosen []Card)
	comb = func(start int, chosen []Card) {
		if len(chosen) == n {
			combination := make([]Card, n)
			copy(combination, chosen)
			result = append(result, combination)
			return
		}
		for i := start; i <= len(cards)-1; i++ {
			comb(i+1, append(chosen, cards[i]))
		}
	}
	comb(0, []Card{})
	return result
}

func countValues(cards []Card) map[int]int {
	valueCount := make(map[int]int)
	for _, card := range cards {
		valueCount[int(card.Value)]++
	}
	return valueCount
}

func isFlush(cards []Card) bool {
	firstSuit := cards[0].Suit
	for _, card := range cards {
		if card.Suit != firstSuit {
			return false
		}
	}
	return true
}

func isStraight(cards []Card) bool {
	values := []int{}
	for _, card := range cards {
		values = append(values, int(card.Value))
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

func getSortedValues(cards []Card) []int {
	values := []int{}
	for _, card := range cards {
		values = append(values, int(card.Value))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	return values
}

func evaluateHand(cards []Card) (HandRanking, []int) {
	valueCount := countValues(cards)

	// Check for flush and straight
	flush := isFlush(cards)
	straight := isStraight(cards)

	// Royal Flush or Straight Flush
	if flush && straight {
		if cards[0].Value == Ten && cards[1].Value == Jack &&
			cards[2].Value == Queen && cards[3].Value == King && cards[4].Value == Ace {
			return RoyalFlush, []int{int(Ace)} // Highest card is Ace in a Royal Flush
		}
		return StraightFlush, []int{int(cards[4].Value)} // Highest card in the straight flush
	}

	// Four of a Kind
	for value, count := range valueCount {
		if count == 4 {
			return FourOfAKind, []int{value}
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
		return FullHouse, []int{three, pair}
	}

	// Flush
	if flush {
		return Flush, getSortedValues(cards)
	}

	// Straight
	if straight {
		return Straight, []int{int(cards[4].Value)}
	}

	// Three of a Kind
	if three != -1 {
		return ThreeOfAKind, []int{three}
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
		return TwoPair, pairs
	}

	// One Pair
	if len(pairs) == 1 {
		return OnePair, pairs
	}

	// High Card
	return HighCard, getSortedValues(cards)
}

func compareTiebreakers(tiebreaker1, tiebreaker2 []int) int {
	for i := range tiebreaker1 {
		if tiebreaker1[i] > tiebreaker2[i] {
			return 1
		} else if tiebreaker1[i] < tiebreaker2[i] {
			return -1
		}
	}
	return 0
}
