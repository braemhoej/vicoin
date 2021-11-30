package node

type Insn int8

const (
	PeerRequest Insn = 0x00
	PeerReply   Insn = 0x01
	ConnAnn     Insn = 0x02
	DissAnn     Insn = 0x03
	Transaction Insn = 0x04
	BlockAnn    Insn = 0x05
)

type Packet struct {
	Instruction Insn
	Data        interface{}
}
