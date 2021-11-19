package network

type Insn int8

const (
	PeerRequest     Insn = 0x00
	PeerReply       Insn = 0x01
	ConnAnnouncment Insn = 0x02
	Transaction     Insn = 0x03
)

type Packet struct {
	Instruction Insn
	Data        interface{}
}
