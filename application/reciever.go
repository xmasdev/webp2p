package application

import (
	"encoding/binary"
	"sync"

	"github.com/webp2p/ui"
)

type Receiver struct {
	buf []byte
	mu  sync.Mutex
}

func NewReceiver() *Receiver {
	return &Receiver{
		buf: make([]byte, 0),
	}
}

func (r *Receiver) Push(data []byte, uiInstance *ui.UI, peerName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.buf = append(r.buf, data...)

	for {
		// Need at least length prefix
		if len(r.buf) < 4 {
			return
		}

		msgLen := binary.BigEndian.Uint32(r.buf[:4])
		totalLen := int(4 + msgLen)

		if len(r.buf) < totalLen {
			return
		}

		msg := string(r.buf[4:totalLen])
		r.buf = r.buf[totalLen:]

		uiInstance.Log("%s: %s", peerName, msg)
	}
}
