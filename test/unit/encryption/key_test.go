package encryption_test

import (
	"encoding/gob"
	"reflect"
	"testing"
	"vicoin/crypto"
)

func TestPublicKeysCanBeBase64StringEncodedDecoded(t *testing.T) {
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

func TestPrivateKeysCanBeBase64StringEncodedDecoded(t *testing.T) {
	gob.Register(crypto.PrivateKey{})
	_, private, _ := crypto.KeyGen(128)
	privateString, _ := private.ToString()
	privateFromString, err := new(crypto.PrivateKey).FromString(privateString)
	if err != nil {
		t.Error("Error when Base64 decoding key")
	}
	if !reflect.DeepEqual(private, privateFromString) {
		t.Error("Error keys are unequal", private, privateFromString)
	}
}

func TestKeysReturnErrorsWhenDecodingFromNonKeyEncoding(t *testing.T) {
	gob.Register(crypto.PrivateKey{})
	privateString := "private.ToString()"
	gob.Register(crypto.PublicKey{})
	publicString := "public.ToString()"
	_, err1 := new(crypto.PublicKey).FromString(publicString)
	_, err2 := new(crypto.PrivateKey).FromString(privateString)
	if err1 == nil {
		t.Error("No error when Base64 decoding public key")
	}
	if err2 == nil {
		t.Error("No error when Base64 decoding private key")
	}
}
