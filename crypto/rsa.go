package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"reflect"
	"vicoin/encoding"
)

var one = big.NewInt(1)

func KeyGen(bits int) (*PublicKey, *PrivateKey, error) {
	var N, E, D, p, q, p1, q1, p1q1 *big.Int
	p, _ = rand.Prime(rand.Reader, bits/2)
	q, _ = rand.Prime(rand.Reader, bits/2)
	p1 = new(big.Int).Sub(p, one)
	q1 = new(big.Int).Sub(q, one)
	p1q1 = new(big.Int).Mul(p1, q1)
	for {
		// Generate E as co-prime to p1g1
		tmpBytes := make([]byte, bits)
		_, err := rand.Read(tmpBytes)
		if err != nil {
			return nil, nil, err
		}
		E = new(big.Int).SetBytes(tmpBytes)
		p1q1gcd := new(big.Int).GCD(nil, nil, p1q1, E)
		if p1q1gcd.Cmp(one) == 0 {
			break // Co-prime found
		}
	}
	N = new(big.Int).Mul(p, q)
	D = big.NewInt(0).ModInverse(E, p1q1)
	return &PublicKey{N, E}, &PrivateKey{N, D}, nil
}

func Encrypt(object interface{}, publicKey *PublicKey) ([]byte, error) {
	serializedObject, err := encoding.Serialize(object)
	if err != nil {
		return nil, err
	}
	message := new(big.Int).SetBytes(serializedObject)
	return new(big.Int).Exp(message, publicKey.E, publicKey.N).Bytes(), nil
}

func Decrypt(bytes []byte, privateKey *PrivateKey) (interface{}, error) {
	cipher := new(big.Int).SetBytes(bytes)
	plain := new(big.Int).Exp(cipher, privateKey.D, privateKey.N).Bytes()
	object, err := encoding.Deserialize(plain)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func Sign(object interface{}, privateKey *PrivateKey) ([]byte, error) {
	serializedObject, err := encoding.Serialize(object)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(serializedObject)
	return Encrypt(hash[:], &PublicKey{privateKey.N, privateKey.D})
}

func Validate(object interface{}, signature []byte, public *PublicKey) (bool, error) {
	decryptedSignature, err1 := Decrypt(signature, &PrivateKey{public.N, public.E})
	if err1 != nil {
		return false, err1
	}
	serializedObject, err2 := encoding.Serialize(object)
	if err2 != nil {
		return false, err2
	}
	hash := sha256.Sum256(serializedObject)
	return reflect.DeepEqual(decryptedSignature, hash[:]), nil
}
