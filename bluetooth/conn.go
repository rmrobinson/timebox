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
func (btc *Connection) Connect(remote *Address, channel int) error {
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
// Read can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetReadDeadline.
func (btc *Connection) Read(b []byte) (n int, err error) {
	return unix.Read(btc.fd, b)
}

// Write writes data to the connection.
// Write can be made to time out and return an error after a fixed
// time limit; see SetDeadline and SetWriteDeadline.
func (btc *Connection) Write(b []byte) (n int, err error) {
	return unix.Write(btc.fd, b)
}

func (btc *Connection) Close() error {
	if err := unix.Close(btc.fd); err != nil {
		return err
	}
	btc.fd = 0
	return nil
}

func (btc *Connection) LocalAddr() net.Addr {
	return btc.local
}

func (btc *Connection) RemoteAddr() net.Addr {
	return btc.remote
}

func (btc *Connection) SetDeadline(t time.Time) error {
	return nil
}

func (btc *Connection) SetReadDeadline(t time.Time) error {
	return nil
}

func (btc *Connection) SetWriteDeadline(t time.Time) error {
	return nil
}
