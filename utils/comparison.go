package utils

import "net"

func SameAddr(a, b *net.UDPAddr) bool {
	if a == nil || b == nil {
		return false
	}

	if a.Port != b.Port {
		return false
	}

	return a.IP.Equal(b.IP)
}
