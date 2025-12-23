package packet

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/webp2p/header"
	"github.com/webp2p/types"
)

func SendData(
	conn *net.UDPConn,
	peer *types.Peer,
	peerMu *sync.Mutex,
	data []byte,
) error {

	peerMu.Lock()
	if !peer.Alive {
		peerMu.Unlock()
		return fmt.Errorf("peer not alive")
	}

	chunks := ChunkData(data, header.MaxPayloadSize)

	for _, chunk := range chunks {
		seq := peer.NextSeq
		peer.NextSeq++

		packet := BuildDataPacket(chunk, seq)

		peer.SendBuffer[seq] = packet
		peer.SendTimes[seq] = time.Now()

		// send immediately once
		conn.WriteToUDP(packet, peer.Addr)
	}

	peerMu.Unlock()
	return nil
}
