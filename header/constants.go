package header

import "time"

const HeaderSize = 12
const HeaderVersion = 1

const KeepAliveInterval = 50 * time.Millisecond
const PeerTimeout = 5 * KeepAliveInterval

const MaxPayloadSize = 12
