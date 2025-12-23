package handlers

import (
	"net"
	"sync"
	"time"

	"github.com/webp2p/application"
	"github.com/webp2p/header"
	"github.com/webp2p/packet"
	"github.com/webp2p/types"
	"github.com/webp2p/ui"
)

func DataPacketHandler(
	conn *net.UDPConn,
	h header.Header,
	payload []byte,
	peer *types.Peer,
	peerMu *sync.Mutex,
	receiver *application.Receiver,
	uiInstance *ui.UI,
) {
	var deliver [][]byte
	var ack uint32

	peerMu.Lock()

	seq := h.SeqNum

	if seq == peer.ExpectedSeq {
		// in-order packet
		deliver = append(deliver, payload)
		peer.ExpectedSeq++

		// drain buffered packets
		for {
			if p, ok := peer.RecvBuffer[peer.ExpectedSeq]; ok {
				deliver = append(deliver, p)
				delete(peer.RecvBuffer, peer.ExpectedSeq)
				peer.ExpectedSeq++
			} else {
				break
			}
		}
	} else if seq > peer.ExpectedSeq {
		// out-of-order, buffer
		if _, exists := peer.RecvBuffer[seq]; !exists {
			peer.RecvBuffer[seq] = payload
		}
	}
	// seq < ExpectedSeq → duplicate, ignore

	ack = peer.ExpectedSeq - 1
	peerMu.Unlock()

	// Deliver ordered data to application
	for _, p := range deliver {
		receiver.Push(p, uiInstance, peer.Name)
	}

	// Send cumulative ACK
	ackPacket := packet.BuildAckPacket(ack)
	conn.WriteToUDP(ackPacket, peer.Addr)
}

func AckPacketHandler(h header.Header, peer *types.Peer, peerMu *sync.Mutex) {
	peerMu.Lock()
	defer peerMu.Unlock()

	ack := h.SeqNum

	// Ignore stale or duplicate ACKs
	if ack <= peer.LastAckedSeq {
		return
	}

	peer.LastAckedSeq = ack

	for seq := range peer.SendBuffer {
		if seq <= ack {
			delete(peer.SendBuffer, seq)
			delete(peer.SendTimes, seq)
		}
	}
}

func handlePacket(conn *net.UDPConn, h header.Header, payload []byte, peer *types.Peer, peerMu *sync.Mutex, uiInstance *ui.UI, receiver *application.Receiver) {
	peerMu.Lock()
	peer.LastSeen = time.Now()
	if !peer.Alive {
		peer.Alive = true
		uiInstance.Log("Peer is now Alive")
	}
	peerMu.Unlock()

	switch h.Type {
	case types.PacketKeepAlive:

	case types.PacketData:
		DataPacketHandler(conn, h, payload, peer, peerMu, receiver, uiInstance)

	case types.PacketAck:
		AckPacketHandler(h, peer, peerMu)
	default:
		uiInstance.Log("Unknown header type from %s: %d", peer.Addr, h.Type)
	}
}
