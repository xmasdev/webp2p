package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/webp2p/application"
	"github.com/webp2p/handlers"
	"github.com/webp2p/types"
	"github.com/webp2p/ui"
	"github.com/webp2p/utils"
)

func main() {
	noStun := flag.Bool("no-stun", false, "Skip STUN public address discovery")

	flag.Parse()
	args := flag.Args()

	// Expected remaining args: <localPort> <me|connect>
	if len(args) != 2 {
		utils.PrintHelp()
		return
	}

	localPort := args[0]
	option := args[1]

	// me is meaningless without stun
	if option == "me" && *noStun {
		fmt.Println("Error: 'me' requires STUN. Remove --no-stun.")
		return
	}

	// Initialize UI
	uiInstance := ui.NewUI()

	// start listening socket
	localAddr, err := net.ResolveUDPAddr("udp", ":"+localPort)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Start UI in a goroutine
	go func() {
		if err := uiInstance.App.Run(); err != nil {
			panic(err)
		}
	}()

	// Give UI time to initialize
	time.Sleep(500 * time.Millisecond)

	uiInstance.Log("Listening to UDP connections on %v", localAddr)

	if !*noStun {
		publicAddr, err := utils.DiscoverPublicAddr(conn)
		if err != nil {
			uiInstance.Log("STUN failed: %v", err)
			return
		}
		uiInstance.Log("My public address: %v", publicAddr)
	} else {
		uiInstance.Log("STUN disabled (--no-stun)")
	}

	// if option was me then return after printing our IP address
	if option == "me" {
		return
	}

	// if option is not connect then return
	if option != "connect" {
		uiInstance.Log("Invalid option: %s", option)
		return
	}

	// Get peer info from UI input
	uiInstance.Log("Enter peer name (type in input box below)")
	peerName := strings.TrimSpace(<-uiInstance.GetInputChannel())
	uiInstance.Log("Peer name set to: %s", peerName)

	uiInstance.Log("Enter peer address")
	peerAddrString := strings.TrimSpace(<-uiInstance.GetInputChannel())
	uiInstance.Log("Peer address set to: %s", peerAddrString)

	peerAddr, err := net.ResolveUDPAddr("udp", peerAddrString)
	if err != nil {
		uiInstance.Log("Error parsing peer address: %v", err)
		panic(err)
	}
	peer := types.NewPeer(peerName, peerAddr)
	var peerMu sync.Mutex

	uiInstance.Log("Trying to connect to peer on %v", peer.Addr)

	// Create Receiver
	receiver := application.NewReceiver()

	// Start handler loops

	go handlers.StartHandlerLoops(conn, peer, &peerMu, uiInstance, receiver)
	go application.InputLoop(conn, peer, &peerMu, uiInstance)

	// Block to prevent exit
	select {}
}
