package types

import (
	"net"
	"time"
)

type Peer struct {
	Name string
	Addr *net.UDPAddr

	// Liveness
	Alive    bool
	LastSeen time.Time

	// --- Sending side ---
	NextSeq      uint32               // next sequence to assign
	LastAckedSeq uint32               // highest contiguous ACK
	SendBuffer   map[uint32][]byte    // seq -> packet bytes
	SendTimes    map[uint32]time.Time // seq -> last send time

	// --- Receiving side ---
	ExpectedSeq uint32            // next expected seq
	RecvBuffer  map[uint32][]byte // out-of-order packets
}

func NewPeer(name string, addr *net.UDPAddr) *Peer {
	return &Peer{
		Name:     name,
		Addr:     addr,
		Alive:    false,
		LastSeen: time.Now(),

		NextSeq:      1,
		LastAckedSeq: 0,
		ExpectedSeq:  1,

		SendBuffer: make(map[uint32][]byte),
		SendTimes:  make(map[uint32]time.Time),
		RecvBuffer: make(map[uint32][]byte),
	}
}

const (
	PacketKeepAlive uint8 = 0x01
	PacketData      uint8 = 0x02
	PacketHandshake uint8 = 0x03
	PacketAck       uint8 = 0x04
)
