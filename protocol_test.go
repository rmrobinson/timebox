package timebox

import (
	"bytes"
	"testing"
)

func TestEncodeCmd(t *testing.T) {
	tests := map[string]struct {
		cmd    byte
		args   []byte
		result []byte
	}{
		"display clock in 12h format": {cmd: CmdSetView, args: []byte{ViewClock, 0x01}, result: []byte{0x05, 0x00, 0x45, 0x00, 0x01}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := getCommandPayload(tc.cmd, tc.args)
			if !bytes.Equal(result, tc.result) {
				t.Fatalf("got %x, expected %x", result, tc.result)
			}
		})
	}
}
