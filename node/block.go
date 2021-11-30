package node

import (
	"vicoin/crypto"
	"vicoin/encoding"
)

type Block struct {
	SequenceNumber int
	Transactions   []string
}

type SignedBlock struct {
	Block     Block
	Signature string
}

func (block *Block) Sign(private *crypto.PrivateKey) (*SignedBlock, error) {
	signature, err := crypto.Sign(block, private)
	if err != nil {
		return nil, err
	}
	encodedSignature, err := encoding.ToB64(signature)
	if err != nil {
		return nil, err
	}
	signedBlock := SignedBlock{
		Block:     *block,
		Signature: encodedSignature,
	}
	return &signedBlock, nil
}

func (block *SignedBlock) Validate(public *crypto.PublicKey) (bool, error) {
	decodedSignature, err := encoding.FromB64(block.Signature)
	if err != nil {
		return false, err
	}
	signature := decodedSignature.([]byte)
	return crypto.Validate(block.Block, signature, public)
}
