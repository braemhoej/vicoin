package crypto

import (
	"errors"
	"math/big"
	"vicoin/encoding"
)

type PrivateKey struct {
	N *big.Int
	D *big.Int
}

func (privateKey *PrivateKey) ToString() (string, error) {
	return encoding.ToB64(privateKey)
}

func (privateKey *PrivateKey) FromString(str string) (*PrivateKey, error) {
	encodedKey, err := encoding.FromB64(str)
	if err != nil {
		return nil, err
	}
	switch key := encodedKey.(type) {
	case PrivateKey:
		privateKey = &key
	default:
		return nil, errors.New("decoded string is not a private key")
	}
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
	encodedKey, err := encoding.FromB64(str)
	if err != nil {
		return nil, err
	}
	switch key := encodedKey.(type) {
	case PublicKey:
		publicKey = &key
	default:
		return nil, errors.New("decoded string is not a public key")
	}
	return publicKey, nil
}
