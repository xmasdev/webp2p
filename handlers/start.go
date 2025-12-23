package handlers

import (
	"net"
	"sync"

	"github.com/webp2p/application"
	"github.com/webp2p/types"
	"github.com/webp2p/ui"
)

func StartHandlerLoops(conn *net.UDPConn, peer *types.Peer, peerMu *sync.Mutex, uiInstance *ui.UI, receiver *application.Receiver) error {

	go ReadLoop(conn, peer, peerMu, uiInstance, receiver)
	go KeepAliveLoop(conn, peer, uiInstance)
	go LivenessMonitorLoop(peer, peerMu, uiInstance)
	go SendLoop(conn, peer, peerMu, uiInstance)
	select {}
}
