# webp2p

A TCP-like protocol implementation over UDP in Go, featuring NAT traversal through UDP hole punching.

## Features

- **TCP-like Protocol over UDP**: Reliable communication layer built on top of UDP
- **NAT Traversal**: Connect to peers across NAT boundaries using UDP hole punching
- **CLI Tools**: Discover your public IP address and establish peer-to-peer connections
- **Terminal UI**: Interactive interface for managing connections
- **Reliable Delivery**: ACKs, retransmission, and packet framing for data integrity
- **Acknowledgement Window**: Sliding window mechanism for efficient packet acknowledgement and flow control


## Getting Started

### Prerequisites
- Go 1.16 or higher

### Installation

```bash
git clone https://github.com/yourusername/webp2p.git
cd webp2p
go build
```


## Architecture

- Custom UDP-based protocol with TCP-like semantics (ACKs, Retransmission)
- UDP hole punching for connection establishment
- Terminal UI for interactive peer management


## License

MIT
