package encryption_test

import (
	"reflect"
	"testing"
	"vicoin/crypto"
)

var lorem128bytes = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas pulvinar ex eget lacus ornare, vel rhoncus velit orci aliquam. "

func TestRSAEncryptDecryptEqualsPlain(t *testing.T) {
	public, private, _ := crypto.KeyGen(2048)
	encMessage, signErr := crypto.Encrypt(lorem128bytes, public)
	decMessage, decErr := crypto.Decrypt(encMessage, private)
	if signErr != nil || decErr != nil {
		t.Error("Error occurred during encryption or decryption")
	}
	if !reflect.DeepEqual(lorem128bytes, decMessage) {
		t.Error("Unexpected result of encryption; ", decMessage, " : expected ", lorem128bytes)
	}
}

func TestRSACorrectlySignedObjectsCanBeValidatedWithProperKey(t *testing.T) {
	public, private, _ := crypto.KeyGen(2048)
	signature, signErr := crypto.Sign(lorem128bytes, private)
	valid, valErr := crypto.Validate(lorem128bytes, signature, public)
	if signErr != nil || valErr != nil {
		t.Error("Error occurred during signing or validation")
	}
	if !valid {
		t.Error("Unable to validate a correctly signed object")
	}
}
func TestRSACorrectlySignedObjectsCantBeValidatedWithImproperKey(t *testing.T) {
	_, private, _ := crypto.KeyGen(2048)
	improperPublic, _, _ := crypto.KeyGen(2048)
	signature, signErr := crypto.Sign(lorem128bytes, private)
	valid, _ := crypto.Validate(lorem128bytes, signature, improperPublic)
	if signErr != nil {
		t.Error("Error occurred during signing")
	}
	if valid {
		t.Error("Validated an correctly signed object with improper key")
	}
}

func TestRSAIncorrectlySignedObjectsCantBeValidated(t *testing.T) {
	_, private, _ := crypto.KeyGen(2048)
	improperPublic, _, _ := crypto.KeyGen(2048)
	signature, signErr := crypto.Sign(lorem128bytes, private)
	valid, _ := crypto.Validate("lorem128bytes", signature, improperPublic)
	if signErr != nil {
		t.Error("Error occurred during signing")
	}
	if valid {
		t.Error("Validated an incorrectly signed object")
	}
}
