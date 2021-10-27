package encryption

import (
	"encoding/gob"
	"reflect"
	"testing"
	"vicoin/crypto"
)

func TestPublicKeysCanBeBase64StringEncoded(t *testing.T) {
	gob.Register(crypto.PublicKey{})
	public, _, _ := crypto.KeyGen(128)
	_, err := public.ToString()
	if err != nil {
		t.Error("Error when Base64 encoding key", err)
	}
}
func TestPublicKeysCanBeBase64StringDecoded(t *testing.T) {
	gob.Register(crypto.PublicKey{})
	public, _, _ := crypto.KeyGen(128)
	publicString, _ := public.ToString()
	publicFromString, err := new(crypto.PublicKey).FromString(publicString)
	if err != nil {
		t.Error("Error when Base64 decoding key")
	}
	if !reflect.DeepEqual(public, publicFromString) {
		t.Error("Error keys are unequal", public, publicFromString)
	}
}
