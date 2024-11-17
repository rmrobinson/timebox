package timebox

import (
	"errors"
)

var (
	ErrMalformedPayload = errors.New("malformed payload received")
	ErrInvalidChecksum  = errors.New("invalid checksum calculated")
	ErrInvalidLength    = errors.New("invalid length received")
)

type message struct {
	payload []byte
}

func newMessage(payload []byte) *message {
	return &message{payload: payload}
}

// encode takes the stored message, calculates the checksum, encodes the message and returns the Timebox-formatted byte array.
func (m *message) encode() []byte {
	var csPayload []byte

	csPayload = append(csPayload, m.payload...)
	csum := checksum(m.payload)
	csPayload = append(csPayload, csum...)

	var escPayload []byte
	escPayload = append(escPayload, 0x01)
	escPayload = append(escPayload, escape(csPayload)...)
	escPayload = append(escPayload, 0x02)

	return escPayload
}

// decode takes a Timebox-formatted byte array, decodes it, validates the checksum and then saves the message.
func (m *message) decode(payload []byte) error {
	if len(payload) < 4 {
		return ErrInvalidLength
	} else if payload[0] != 0x01 || payload[len(payload)-1] != 0x02 {
		return ErrMalformedPayload
	}

	unescaped, err := unescape(payload[1 : len(payload)-1])
	if err != nil {
		return err
	}

	csum := checksum(unescaped[0 : len(unescaped)-2])
	if csum[0] != unescaped[len(unescaped)-2] || csum[1] != unescaped[len(unescaped)-1] {
		return ErrInvalidChecksum
	}

	m.payload = unescaped[:len(unescaped)-2]
	return nil
}
