package ui

import "image"

// This package contains the layout of the game, including the size of the screen, the size of the map, and the size of the message log.

const Y_AXIS_AJUSTMENT = 3

// Table position
const (
	TABLE_CENTER_X = 65
	TABLE_CENTER_Y = 20 + Y_AXIS_AJUSTMENT
	TABLE_RADIUS_X = 55
	TABLE_RADIUS_Y = 17
)

// Community cards position
const (
	COMMUNITY_CARDS_X = 45
	COMMUNITY_CARDS_Y = 15 + Y_AXIS_AJUSTMENT
)

// Pot position
const (
	POT_X = 60
	POT_Y = 24 + Y_AXIS_AJUSTMENT
)

// Pocket pair position
const (
	POCKET_PAIR_X = 54
	POCKET_PAIR_Y = 34 + Y_AXIS_AJUSTMENT
)

// Control panel position
const (
	CONTROL_PANEL_X_LEFT  = 12
	CONTROL_PANEL_X_RIGHT = 120
	CONTROL_PANEL_Y       = 45 + Y_AXIS_AJUSTMENT
)

// Balance info list position
const (
	BALANCE_INFO_X = 150
	BALANCE_INFO_Y = 2 + Y_AXIS_AJUSTMENT
)

// Define type
type Layout map[int]image.Point

// map to store other players' positions in layout
var PLAYER_LAYOUT = map[int]map[int]image.Point{
	// 1 other player
	2: {
		0: {X: POCKET_PAIR_X, Y: POCKET_PAIR_Y},
		1: {X: 59, Y: 1 + Y_AXIS_AJUSTMENT},
	},
	// 2 other players
	3: {
		0: {X: POCKET_PAIR_X, Y: POCKET_PAIR_Y},
		1: {X: 23, Y: 5 + Y_AXIS_AJUSTMENT},
		2: {X: 94, Y: 5 + Y_AXIS_AJUSTMENT},
	},
	// 3 other players
	4: {
		0: {X: POCKET_PAIR_X, Y: POCKET_PAIR_Y},
		1: {X: 5, Y: 17 + Y_AXIS_AJUSTMENT},
		2: {X: 59, Y: 1 + Y_AXIS_AJUSTMENT},
		3: {X: 112, Y: 17 + Y_AXIS_AJUSTMENT}},
	// 4 other players
	5: {
		0: {X: POCKET_PAIR_X, Y: POCKET_PAIR_Y},
		1: {X: 8, Y: 25 + Y_AXIS_AJUSTMENT},
		2: {X: 27, Y: 4 + Y_AXIS_AJUSTMENT},
		3: {X: 88, Y: 4 + Y_AXIS_AJUSTMENT},
		4: {X: 109, Y: 25 + Y_AXIS_AJUSTMENT}},
	// 5 other players
	6: {
		0: {X: POCKET_PAIR_X, Y: POCKET_PAIR_Y},
		1: {X: 15, Y: 26 + Y_AXIS_AJUSTMENT},
		2: {X: 17, Y: 7 + Y_AXIS_AJUSTMENT},
		3: {X: 59, Y: 1 + Y_AXIS_AJUSTMENT},
		4: {X: 100, Y: 7 + Y_AXIS_AJUSTMENT},
		5: {X: 102, Y: 26 + Y_AXIS_AJUSTMENT}},
}
