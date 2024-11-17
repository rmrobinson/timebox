package bluetooth

import (
	"errors"
	"strconv"
	"strings"
)

// Address stores a parsed address for use in establishing a Bluetooth connection
type Address struct {
	str   string
	bytes [6]byte
}

// NewAddress takes the supplied Bluetooth MAC and converts it into an address for use.
// If the string isn't a validly formatted MAC, it will error
func NewAddress(addr string) (*Address, error) {
	bta := &Address{
		str: addr,
	}

	if err := bta.parse(); err != nil {
		return nil, err
	}

	return bta, nil
}

// Network returns the "bt" network for conforming to the net.Addr interface
func (bta *Address) Network() string {
	return "bt"
}

// String returns the formatted MAC address for conforming to the net.Addr interface
func (bta *Address) String() string {
	return bta.str
}

func (bta *Address) parse() error {
	splitAddr := strings.Split(bta.str, ":")
	if len(splitAddr) != 6 {
		return errors.New("invalid address format")
	}

	for i, addrByte := range splitAddr {
		b, err := strconv.ParseUint(addrByte, 16, 8)
		if err != nil {
			return err
		}
		bta.bytes[5-i] = byte(b)
	}

	return nil
}
