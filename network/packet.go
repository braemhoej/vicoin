package network

type Insn int8

const (
	PeerRequest     Insn = 0
	PeerReply       Insn = 1
	ConnAnnouncment Insn = 2
	Transaction     Insn = 3
)

type Packet struct {
	Instruction Insn
	Data        interface{}
}
