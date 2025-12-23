package utils

import (
	"net"

	"github.com/pion/stun"
)

func DiscoverPublicAddr(conn *net.UDPConn) (*net.UDPAddr, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", "stun.l.google.com:19302")
	if err != nil {
		return nil, err
	}

	message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	_, err = conn.WriteToUDP(message.Raw, serverAddr)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1500)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}

	var res stun.Message
	res.Raw = buf[:n]
	if err := res.Decode(); err != nil {
		return nil, err
	}

	var xorAddr stun.XORMappedAddress
	if err := xorAddr.GetFrom(&res); err != nil {
		return nil, err
	}

	return &net.UDPAddr{
		IP:   xorAddr.IP,
		Port: xorAddr.Port,
	}, nil
}
