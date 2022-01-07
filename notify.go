package diameter

import "log"

type Direction bool

func (v Direction) String() string {
	if v {
		return "Tx"
	}
	return "Rx"
}

const (
	Tx Direction = true
	Rx Direction = false
)

var TraceMessage = func(msg Message, dct Direction, err error) {
	log.Println(dct, "message handling: error=", err, "\n", msg)
}

var TraceState = func(old, new, event string, err error) {
	log.Println("state update:", old, "->", new, "by event", event, "error", err)
}
