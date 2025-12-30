//go:build !linux

package bluetooth

import "errors"

// Connect returns that this isn't supported on this platform
func (btc *Connection) Connect(remote *Address, channel uint8) error {
	return errors.New("not supported")
}

// Read reads data from the connection.
func (btc *Connection) Read(b []byte) (n int, err error) {
	return -1, errors.New("not supported")
}

// Write writes data to the connection.
func (btc *Connection) Write(b []byte) (n int, err error) {
	return -1, errors.New("not supported")
}

// Close closes the connection.
func (btc *Connection) Close() error {
	return nil
}
