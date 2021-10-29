package encoding

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

func Serialize(object interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(&object)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func Deserialize(data []byte) (interface{}, error) {
	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	var object interface{}
	err := dec.Decode(&object)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func ToB64(object interface{}) (string, error) {
	serializedObject, err := Serialize(object)
	if err != nil {
		return "", err
	}
	b64String := base64.StdEncoding.EncodeToString(serializedObject)
	return b64String, nil
}
func FromB64(b64String string) (interface{}, error) {
	serializedObject, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		return nil, err
	}
	object, _ := Deserialize(serializedObject)
	return object, nil
}
