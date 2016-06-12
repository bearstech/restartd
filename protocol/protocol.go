package protocol

import "errors"

type Command int

const (
	Status Command = iota
	Start
	Stop
	Restart
	Reload
	Unknown
)

var Commands = []string{"status", "start", "stop", "restart", "reload"}

func (c Command) Command() string {
	return Commands[int(c)]
}

type Message struct {
	Service string
	Command Command
}

type Response struct {
	Code    int
	Message string
}

func ParseCommand(txt string) (Command, error) {
	for i := 0; i < len(Commands); i++ {
		if Commands[i] == txt {
			return Command(i), nil
		}
	}
	return Unknown, errors.New("Command unknown")
}
