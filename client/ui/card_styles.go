// Description: This file contains styles for the cards.
// The cardStyle map contains the style for each suit.

package ui

import (
	msgpb "go-pk-server/gen"
)

var cardStyle = map[msgpb.SuitType]Style{
	msgpb.SuitType_HEARTS:   {ColorRed, ColorBlack, ModifierBold},
	msgpb.SuitType_DIAMONDS: {ColorYellow, ColorBlack, ModifierBold},
	msgpb.SuitType_CLUBS:    {ColorBlue, ColorBlack, ModifierBold},
	msgpb.SuitType_SPADES:   {ColorWhite, ColorBlack, ModifierBold},
}

var suitsIcon = map[msgpb.SuitType]rune{
	msgpb.SuitType_HEARTS:   '♥',
	msgpb.SuitType_DIAMONDS: '♦',
	msgpb.SuitType_CLUBS:    '♣',
	msgpb.SuitType_SPADES:   '♠',
}

var ranksIcon = map[msgpb.RankType]rune{
	msgpb.RankType_ACE:   'A',
	msgpb.RankType_DEUCE: '2',
	msgpb.RankType_THREE: '3',
	msgpb.RankType_FOUR:  '4',
	msgpb.RankType_FIVE:  '5',
	msgpb.RankType_SIX:   '6',
	msgpb.RankType_SEVEN: '7',
	msgpb.RankType_EIGHT: '8',
	msgpb.RankType_NINE:  '9',
	msgpb.RankType_TEN:   'T',
	msgpb.RankType_JACK:  'J',
	msgpb.RankType_QUEEN: 'Q',
	msgpb.RankType_KING:  'K',
}

var emptyCard = []string{
	"┌─────┐",
	"│.....│",
	"│.....│",
	"│.....│",
	"└─────┘",
}
