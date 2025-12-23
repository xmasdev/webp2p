package packet

import (
	"github.com/webp2p/header"
	"github.com/webp2p/types"
)

func BuildDataPacket(payload []byte, seq uint32) []byte {
	h := header.Header{
		Version: header.HeaderVersion,
		Type:    types.PacketData,
		Flags:   0,
		SeqNum:  seq,
		Length:  uint32(len(payload)),
	}.EncodeHeader()

	packet := make([]byte, header.HeaderSize+len(payload))
	copy(packet[:header.HeaderSize], h)
	copy(packet[header.HeaderSize:], payload)
	return packet
}

func BuildAckPacket(ack uint32) []byte {
	h := header.Header{
		Version: header.HeaderVersion,
		Type:    types.PacketAck,
		Flags:   0,
		SeqNum:  ack,
		Length:  0,
	}.EncodeHeader()
	packet := make([]byte, header.HeaderSize)
	copy(packet[:header.HeaderSize], h)
	return packet
}

func ChunkData(data []byte, maxPayload int) [][]byte {
	var chunks [][]byte
	for len(data) > 0 {
		n := min(len(data), maxPayload)
		chunks = append(chunks, data[:n])
		data = data[n:]
	}
	return chunks
}
