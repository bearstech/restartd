package model

import "errors"

func ParseCommand(txt string) (Message_Commands, error) {
	v, ok := Message_Commands_value[txt]
	if ok {
		return Message_Commands(v), nil
	} else {
		return Message_Commands(-1), errors.New("Command unknown")
	}
}
