package ui

import "github.com/nsf/termbox-go"

type KeyboardEventType int

// Keyboard events
const (
	ENTER     KeyboardEventType = 0
	EXIT      KeyboardEventType = 2
	LEFT      KeyboardEventType = 3
	RIGHT     KeyboardEventType = 4
	UP        KeyboardEventType = 5
	DOWN      KeyboardEventType = 6
	SPACE     KeyboardEventType = 7
	BACKSPACE KeyboardEventType = 8
	MENU1     KeyboardEventType = 9
	MENU2     KeyboardEventType = 10
	MENU3     KeyboardEventType = 11

	START_GAME KeyboardEventType = 20
	PAUSE_GAME KeyboardEventType = 21
	LEAVE_GAME KeyboardEventType = 22
	TAKE_BUYIN KeyboardEventType = 23
	GIVE_BUYIN KeyboardEventType = 24

	FOLD          KeyboardEventType = 30
	CHECK_OR_CALL KeyboardEventType = 31
	RAISE         KeyboardEventType = 32
	ALL_IN        KeyboardEventType = 33
)

type KeyboardEvent struct {
	EventType KeyboardEventType
	Key       termbox.Key
}

// func keyToDirection(k termbox.Key) direction {
// 	switch k {
// 	case termbox.KeyArrowLeft:
// 		return LEFT
// 	case termbox.KeyArrowDown:
// 		return DOWN
// 	case termbox.KeyArrowRight:
// 		return RIGHT
// 	case termbox.KeyArrowUp:
// 		return UP
// 	default:
// 		return 0
// 	}
// }

func ListenToKeyboard(evChan chan KeyboardEvent) {
	termbox.SetInputMode(termbox.InputEsc)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowLeft:
				evChan <- KeyboardEvent{EventType: LEFT, Key: ev.Key}
			case termbox.KeyArrowDown:
				evChan <- KeyboardEvent{EventType: DOWN, Key: ev.Key}
			case termbox.KeyArrowRight:
				evChan <- KeyboardEvent{EventType: RIGHT, Key: ev.Key}
			case termbox.KeyArrowUp:
				evChan <- KeyboardEvent{EventType: UP, Key: ev.Key}
			case termbox.KeyEnter:
				evChan <- KeyboardEvent{EventType: ENTER, Key: ev.Key}
			case termbox.KeySpace:
				evChan <- KeyboardEvent{EventType: SPACE, Key: ev.Key}
			case termbox.KeyBackspace2:
				evChan <- KeyboardEvent{EventType: BACKSPACE, Key: ev.Key}
			case termbox.KeyEsc:
				evChan <- KeyboardEvent{EventType: EXIT, Key: ev.Key}
			default:
				switch ev.Ch {
				case 's':
					evChan <- KeyboardEvent{EventType: START_GAME, Key: ev.Key}
				case 'p':
					evChan <- KeyboardEvent{EventType: PAUSE_GAME, Key: ev.Key}
				case 'l':
					evChan <- KeyboardEvent{EventType: LEAVE_GAME, Key: ev.Key}
				case 't':
					evChan <- KeyboardEvent{EventType: TAKE_BUYIN, Key: ev.Key}
				case 'g':
					evChan <- KeyboardEvent{EventType: GIVE_BUYIN, Key: ev.Key}
				case 'f':
					evChan <- KeyboardEvent{EventType: FOLD, Key: ev.Key}
				case 'c':
					evChan <- KeyboardEvent{EventType: CHECK_OR_CALL, Key: ev.Key}
				case 'r':
					evChan <- KeyboardEvent{EventType: RAISE, Key: ev.Key}
				case 'a':
					evChan <- KeyboardEvent{EventType: ALL_IN, Key: ev.Key}
				case '1':
					evChan <- KeyboardEvent{EventType: MENU1, Key: ev.Key}
				case '2':
					evChan <- KeyboardEvent{EventType: MENU2, Key: ev.Key}
				case '3':
					evChan <- KeyboardEvent{EventType: MENU3, Key: ev.Key}
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
