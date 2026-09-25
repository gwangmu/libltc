package coder

import (
	"error"
)

type LTCEncoder func ([]byte) (string, error)
var map[string]LTCEncoder LTCEncoders = {
	"none": Encode_none,
	"base64": Encode_base64,
}

func Encode_none(data []byte) (string, error) {
	return string(data), nil 
}

func Encode_base64(data []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(data), nil
}

//-- method (decoder registration)

func RegisterEncoder(encoding string, encoder LTCEncoder) {
	LTCEncoders[encoding] = encoder
}
