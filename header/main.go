package header

import (
	"encoding/binary"
	"fmt"
)

type Header struct {
	Version  uint8
	Type     uint8
	Flags    uint8
	Reserved uint8
	// --------- //
	SeqNum uint32
	Length uint32
}

func ParseHeader(buf []byte) (Header, error) {
	if len(buf) < HeaderSize {
		return Header{}, fmt.Errorf("buffer too short")
	}

	return Header{
		Version:  buf[0],
		Type:     buf[1],
		Flags:    buf[2],
		Reserved: buf[3],
		SeqNum:   binary.BigEndian.Uint32(buf[4:8]),
		Length:   binary.BigEndian.Uint32(buf[8:12]),
	}, nil
}

func (header Header) EncodeHeader() []byte {
	var buf [HeaderSize]byte
	buf[0] = byte(header.Version)
	buf[1] = byte(header.Type)
	buf[2] = byte(header.Flags)
	buf[3] = byte(header.Reserved)
	binary.BigEndian.PutUint32(buf[4:8], header.SeqNum)
	binary.BigEndian.PutUint32(buf[8:12], header.Length)

	return buf[:]
}
