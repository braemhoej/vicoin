package registration

import (
	"encoding/gob"
	"vicoin/account"
	"vicoin/crypto"
	"vicoin/node"
)

func RegisterStructsWithGob() {
	gob.Register(account.SignedTransaction{})
	gob.Register(account.Transaction{})
	gob.Register(crypto.PrivateKey{})
	gob.Register(crypto.PublicKey{})
	gob.Register(node.Block{})
}
