package node

type Insn int8

const (
	PeerRequest Insn = 0x00
	PeerReply   Insn = 0x01
	ConnAnn     Insn = 0x02
	DissAnn     Insn = 0x03
	Heartbeat   Insn = 0x04
	Transaction Insn = 0x05
)

type Packet struct {
	Instruction Insn
	Data        interface{}
	VectorClock VectorClock
}

type VectorClock struct {
	vector map[Peer]int
}

func NewVectorClock() VectorClock {
	return VectorClock{
		vector: make(map[Peer]int),
	}
}

func (vectorClock *VectorClock) Get(peer Peer) int {
	return vectorClock.vector[peer]
}

func (vectorClock *VectorClock) Increment(peer Peer) {
	vectorClock.vector[peer] += 1
}

func (vectorClock *VectorClock) Delete(peer Peer) {
	delete(vectorClock.vector, peer)
}
