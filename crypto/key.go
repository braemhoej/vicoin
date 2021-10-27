package crypto

import (
	"encoding/base64"
	"math/big"
)

type PrivateKey struct {
	N *big.Int
	D *big.Int
}
type PublicKey struct {
	N *big.Int
	E *big.Int
}

func (publicKey *PublicKey) ToString() (string, error) {
	serializedKey, err := serialize(publicKey)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(serializedKey)
	return str, nil
}
func (publicKey *PublicKey) FromString(str string) (*PublicKey, error) {
	serializedKey, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	key, _ := deserialize(serializedKey)
	publicKey.E = key.(PublicKey).E
	publicKey.N = key.(PublicKey).N
	return publicKey, nil
}
