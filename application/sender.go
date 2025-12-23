package application

import (
	"encoding/binary"
	"net"
	"sync"

	"github.com/webp2p/packet"
	"github.com/webp2p/types"
	"github.com/webp2p/ui"
)

func InputLoop(conn *net.UDPConn, peer *types.Peer, peerMu *sync.Mutex, uiInstance *ui.UI) {
	for line := range uiInstance.GetInputChannel() {
		uiInstance.Log(">>> %s", line)
		SendMessage(line, conn, peer, peerMu)
	}
}

func SendMessage(text string, conn *net.UDPConn, peer *types.Peer, peerMu *sync.Mutex) {
	buf := make([]byte, 4+len(text))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(text)))
	copy(buf[4:], []byte(text))
	packet.SendData(conn, peer, peerMu, buf)
}
