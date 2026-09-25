package coder

import (
	"error"
	"encoding/base64"
)

type LTCDecoder func (string) ([]byte, error)
var map[string]LTCDecoder LTCDecoders = {
	"none": Decode_none,
	"base64": Decode_base64,
}

func Decode_none(data string) ([]byte, error) {
	return []byte(data), nil
}

func Decode_base64(data string) ([]byte, error) {
	dec, err := base64.StdEncoding.DecodeString(data)
	return dec, err
}

//-- method (decoder registration)

func RegisterDecoder(encoding string, decoder LTCDecoder) {
	LTCDecoders[encoding] = decoder
}
