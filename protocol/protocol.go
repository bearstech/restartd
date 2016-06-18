package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
)

func ParseCommand(txt string) (Message_Commands, error) {
	v, ok := Message_Commands_value[txt]
	if ok {
		return Message_Commands(v), nil
	} else {
		return Message_Commands(-1), errors.New("Command unknown")
	}
}

func Write(wire io.Writer, msg proto.Message) error {
	txt, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = binary.Write(wire, binary.LittleEndian, uint16(len(txt)))
	if err != nil {
		return err
	}
	_, err = wire.Write([]byte(txt))
	if err != nil {
		return err
	}
	return nil
}

func Read(wire io.Reader, msg proto.Message) error {
	var size uint16
	err := binary.Read(wire, binary.LittleEndian, &size)
	if err != nil {
		return err
	}
	buf := make([]byte, size)
	_, err = io.ReadFull(wire, buf)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(buf, msg)
	if err != nil {
		return err
	}
	return nil
}
