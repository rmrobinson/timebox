package bluetooth

import (
	"errors"
	"log"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

// Connection is a minimal wrapper around a UNIX Bluetooth socket; it conforms to the net.Conn interface for use by higher level libraries.
type Connection struct {
	fd int

	local  *Address
	remote *Address
}

// Connect establishes an RFCOMM socket to the specified Bluetooth address using the specified channel.
func (btc *Connection) Connect(remote *Address, channel uint8) error {
	if btc.fd > 0 {
		return errors.New("already connected")
	} else if remote == nil {
		return errors.New("address required")
	}

	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_STREAM, unix.BTPROTO_RFCOMM)
	if err != nil {
		log.Printf("unable to create socket: %s\n", err.Error())
		return err
	}

	if err := unix.Connect(fd, &unix.SockaddrRFCOMM{Addr: remote.bytes, Channel: channel}); err != nil {
		log.Printf("unable to connect: %s\n", err.Error())
		return err
	}

	btc.fd = fd
	btc.remote = remote
	return nil
}

// Read reads data from the connection.
func (btc *Connection) Read(b []byte) (n int, err error) {
	return unix.Read(btc.fd, b)
}

// Write writes data to the connection.
func (btc *Connection) Write(b []byte) (n int, err error) {
	return unix.Write(btc.fd, b)
}

// Close closes the connection.
func (btc *Connection) Close() error {
	if err := unix.Close(btc.fd); err != nil {
		return err
	}
	btc.fd = 0
	return nil
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
