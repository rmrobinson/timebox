package bluetooth

import (
	"net"
	"time"
)

// Connection is a minimal wrapper around a UNIX Bluetooth socket; it conforms to the net.Conn interface for use by higher level libraries.
type Connection struct {
	fd int

	local  *Address
	remote *Address
}

// LocalAddr returns the local address of the connection.
func (btc *Connection) LocalAddr() net.Addr {
	return btc.local
}

// RemoteAddr returns the remote address of the connection.
func (btc *Connection) RemoteAddr() net.Addr {
	return btc.remote
}

// SetDeadline has no effect
func (btc *Connection) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline has no effect
func (btc *Connection) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline has no effect
func (btc *Connection) SetWriteDeadline(t time.Time) error {
	return nil
}
