package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
)

func Write(wire io.Writer, msg proto.Message) error {
	txt, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	size := len(txt)
	if size >= 65536 {
		return errors.New(fmt.Sprintf("Message is too big : %i >= 65536", size))
	}
	err = binary.Write(wire, binary.LittleEndian, uint16(size))
	if err != nil {
		return err
	}
	s, err := wire.Write([]byte(txt))
	if err != nil {
		return err
	}
	if s < size {
		return errors.New("Partial write")
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
	s, err := io.ReadFull(wire, buf)
	if err != nil {
		return err
	}
	if uint16(s) < size {
		return errors.New("Partial read")
	}
	err = proto.Unmarshal(buf, msg)
	if err != nil {
		return err
	}
	return nil
}
