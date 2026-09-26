package coder

import (
	"error"
	"unicode/utf8"
)

type Encoder func ([]byte) (string, error)
var map[string]Encoder Encoders = {
	"none": Encode_none,
	"base64": Encode_base64,
}

func Encode_none(data []byte) (string, error) {
	if !utf8.Valid(data) {
		return "", errors.New("Invalid UTF-8-encoded runes")
	} else {
		return string(data), nil 
	}
}

func Encode_base64(data []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(data), nil
}

//-- method (decoder registration)

func RegisterEncoder(encoding string, encoder Encoder) {
	Encoders[encoding] = encoder
}
