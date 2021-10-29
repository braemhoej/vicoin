package node

import "net"

type Peer struct {
	Addr *net.TCPAddr
}
