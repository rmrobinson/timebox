//go:build linux

package bluetooth

import (
	"errors"
	"log"

	"golang.org/x/sys/unix"
)

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
