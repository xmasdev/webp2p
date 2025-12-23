package header

import "github.com/webp2p/types"

func BuildKeepAliveHeader(seq uint32) []byte {
	return Header{
		Version:  HeaderVersion,
		Type:     types.PacketKeepAlive,
		Flags:    0,
		Reserved: 0,
		SeqNum:   seq,
		Length:   0,
	}.EncodeHeader()
}
