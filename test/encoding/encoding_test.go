package encoding

import (
	"encoding/gob"
	"reflect"
	"testing"
	"vicoin/encoding"
)

func TestObjectSerializationDeserializationEqualsOriginal(t *testing.T) {
	type Object struct {
		Field string
	}
	gob.Register(Object{})
	object := Object{
		Field: "something",
	}
	serializedObject, _ := encoding.Serialize(object)
	deserializedObject, _ := encoding.Deserialize(serializedObject)
	if !reflect.DeepEqual(object, deserializedObject.(Object)) {
		t.Error("Deserialized object doesn't reflect original", deserializedObject, object)
	}
}
