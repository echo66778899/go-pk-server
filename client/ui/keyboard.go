package ui

import "github.com/nsf/termbox-go"

type KeyboardEventType int

// Keyboard events
const (
	ENTER     KeyboardEventType = 0
	REQCHIP   KeyboardEventType = 1
	END       KeyboardEventType = 2
	LEFT      KeyboardEventType = 3
	RIGHT     KeyboardEventType = 4
	UP        KeyboardEventType = 5
	DOWN      KeyboardEventType = 6
	SPACE     KeyboardEventType = 7
	BACKSPACE KeyboardEventType = 8
	MENU1     KeyboardEventType = 9
	MENU2     KeyboardEventType = 10
	MENU3     KeyboardEventType = 11
	NEXT      KeyboardEventType = 12
	PAYBACK   KeyboardEventType = 13
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
				evChan <- KeyboardEvent{EventType: END, Key: ev.Key}
			default:
				switch ev.Ch {
				case 'r':
					evChan <- KeyboardEvent{EventType: REQCHIP, Key: ev.Key}
				case 'p':
					evChan <- KeyboardEvent{EventType: PAYBACK, Key: ev.Key}
				case 'n':
					evChan <- KeyboardEvent{EventType: NEXT, Key: ev.Key}
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
