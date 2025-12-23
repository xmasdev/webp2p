package handlers

import (
	"net"
	"sync"
	"time"

	"github.com/webp2p/application"
	"github.com/webp2p/header"
	"github.com/webp2p/types"
	"github.com/webp2p/ui"
	"github.com/webp2p/utils"
)

func KeepAliveLoop(conn *net.UDPConn, peer *types.Peer, uiInstance *ui.UI) {
	ticker := time.NewTicker(header.KeepAliveInterval)
	defer ticker.Stop()

	var seq uint32 = 0
	for {
		<-ticker.C
		h := header.BuildKeepAliveHeader(seq)
		seq++
		_, err := conn.WriteToUDP(h, peer.Addr)
		if err != nil {
			uiInstance.Log("Error while sending Keep Alive: %v", err)
			continue
		}
	}
}

func ReadLoop(conn *net.UDPConn, peer *types.Peer, peerMu *sync.Mutex, uiInstance *ui.UI, receiver *application.Receiver) {
	packet := make([]byte, 1500)
	for {
		n, addr, err := conn.ReadFromUDP(packet)
		if err != nil {
			uiInstance.Log("Error while reading: %v", err)
			continue
		}
		// handle packet
		if !utils.SameAddr(addr, peer.Addr) {
			uiInstance.Log("Unknown packet from %s | Peer address is %s", addr, peer.Addr)
			continue
		}
		if n < header.HeaderSize {
			uiInstance.Log("Header too small from: %s", addr)
			continue
		}

		h, err := header.ParseHeader(packet[:n])
		if err != nil {
			uiInstance.Log("Error parsing header from %s: %v", addr, err)
			continue
		}
		if h.Version != header.HeaderVersion {
			uiInstance.Log("Invalid protocol version from %s", addr)
			continue
		}

		payloadLen := int(h.Length)
		actualPayloadLen := n - header.HeaderSize

		if payloadLen != actualPayloadLen {
			uiInstance.Log(
				"Length mismatch from %s: header=%d actual=%d",
				addr, payloadLen, actualPayloadLen,
			)
			continue
		}

		payload := packet[header.HeaderSize:n]
		handlePacket(conn, h, payload, peer, peerMu, uiInstance, receiver)
	}

}

func LivenessMonitorLoop(peer *types.Peer, peerMu *sync.Mutex, uiInstance *ui.UI) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		if peer.Alive && time.Since(peer.LastSeen) > header.PeerTimeout {
			peerMu.Lock()
			peer.Alive = false
			uiInstance.Log("Peer is dead")
			peerMu.Unlock()
		}
	}

}

func SendLoop(conn *net.UDPConn, peer *types.Peer, peerMu *sync.Mutex, uiInstance *ui.UI) {
	maxResends := 5
	ticker := time.NewTicker(200 * time.Millisecond)
	for range ticker.C {
		count := 5
		now := time.Now()

		peerMu.Lock()

		if !peer.Alive {
			peerMu.Unlock()
			continue
		}

		for seq, pkt := range peer.SendBuffer {
			if count >= maxResends {
				break
			}
			// Skip packets already ACKed (defensive)
			if seq <= peer.LastAckedSeq {
				continue
			}

			lastSent := peer.SendTimes[seq]
			if now.Sub(lastSent) >= time.Second {
				// retransmit
				_, err := conn.WriteToUDP(pkt, peer.Addr)
				if err == nil {
					peer.SendTimes[seq] = now
				}
			}
			count++
		}

		peerMu.Unlock()
	}
}
