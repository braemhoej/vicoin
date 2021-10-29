package crypto

import (
	"math/big"
	"vicoin/internal/encoding"
)

type PrivateKey struct {
	N *big.Int
	D *big.Int
}

func (privateKey *PrivateKey) ToString() (string, error) {
	return encoding.ToB64(privateKey)
}
func (privateKey *PrivateKey) FromString(str string) (*PrivateKey, error) {
	key, err := encoding.FromB64(str)
	if err != nil {
		return nil, err
	}
	privateKey.D = key.(PrivateKey).D
	privateKey.N = key.(PrivateKey).N
	return privateKey, nil
}

type PublicKey struct {
	N *big.Int
	E *big.Int
}

func (publicKey *PublicKey) ToString() (string, error) {
	return encoding.ToB64(publicKey)
}
func (publicKey *PublicKey) FromString(str string) (*PublicKey, error) {
	key, err := encoding.FromB64(str)
	if err != nil {
		return nil, err
	}
	publicKey.E = key.(PublicKey).E
	publicKey.N = key.(PublicKey).N
	return publicKey, nil
}
