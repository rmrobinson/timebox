package timebox

import (
	"bytes"
	"testing"
)

func TestEncodeMessage(t *testing.T) {
	tests := map[string]struct {
		payload []byte
		result  []byte
	}{
		"display clock in 12h format": {payload: []byte{0x05, 0x00, 0x45, 0x00, 0x01}, result: []byte{0x01, 0x05, 0x00, 0x45, 0x00, 0x03, 0x04, 0x4B, 0x00, 0x02}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := newMessage(tc.payload)
			result := m.encode()
			if !bytes.Equal(result, tc.result) {
				t.Fatalf("got %x, expected %x", result, tc.result)
			}
		})
	}
}
