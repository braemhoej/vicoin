package registration

import (
	"encoding/gob"
	"vicoin/crypto"
	"vicoin/internal/account"
)

func RegisterStructsWithGob() {
	gob.Register(account.SignedTransaction{})
	gob.Register(account.Transaction{})
	gob.Register(crypto.PrivateKey{})
	gob.Register(crypto.PublicKey{})
}
