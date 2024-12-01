package timebox

import (
	"errors"
	"log"
	"net"
	"time"
)

// TemperatureUnit describes the supported units for temperature
type TemperatureUnit int

const (
	Celsius TemperatureUnit = iota
	Farenheight
)

// WeatherCondition describes the current weather conditions
type WeatherCondition int

type Colour struct {
	R byte
	G byte
	B byte
}

// Conn is a connection to a Timebox device
type Conn struct {
	c net.Conn

	colour *Colour
}

// NewConn creates a connection using the existing, established network connection.
// This is likely to be a bluetooth RFCOMM socket.
func NewConn(c net.Conn) *Conn {
	return &Conn{
		c: c,
	}
}

// SetColor sets the colour of this device for subsequent commands.
func (c *Conn) SetColor(colour *Colour) {
	c.colour = colour
}

// Initialize confirms that the connection is properly established to a supported Timebox device.
func (c *Conn) Initialize() error {
	resp := make([]byte, 1024)
	n, err := c.c.Read(resp)
	if err != nil {
		log.Printf("unable to read message: %s\n", err.Error())
		return err
	} else if n < len(helloResponse) {
		log.Printf("received message shorter than hello response\n")
		return errors.New("received payload too short")
	}

	for i, b := range helloResponse {
		if resp[i] != b {
			log.Printf("initial response wasn't as expected: got %x but expected %x\n", resp, helloResponse)
			return errors.New("initial response invalid")
		}
	}

	return nil
}

func (c *Conn) receiveMessage() (*message, error) {
	b := make([]byte, 1024)
	n, err := c.c.Read(b)
	if err != nil {
		log.Printf("unable to read message: %s\n", err.Error())
		return nil, err
	}

	m := &message{}
	if err := m.decode(b[:n]); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Conn) sendMessage(m *message) error {
	payload := m.encode()

	log.Printf("%x\n", payload)

	n, err := c.c.Write(payload)
	if err != nil {
		log.Printf("unable to send message: %s\n", err.Error())
		return err
	} else if n != len(payload) {
		log.Printf("only sent %d but needed to send %d\n", n, len(payload))
		return errors.New("partial send")
	}
	return nil
}

// SetTime sets the time on the Timebox
func (c *Conn) SetTime(t time.Time) error {
	// TODO: validate socket is connected

	encodedTime := []byte{
		byte(t.Year() >> 8),
		byte(t.Year() & 0xff),
		byte(t.Month()),
		byte(t.Day()),
		byte(t.Hour()),
		byte(t.Minute()),
		byte(t.Second()),
		byte(0),
	}

	p := getCommandPayload(CmdSetDateTime, encodedTime)
	return c.sendMessage(newMessage(p))
}

// SetTemperatureAndWeather sets the time and weather on the Timebox
func (c *Conn) SetTemperatureAndWeather(temp int, tempUnit TemperatureUnit, weather WeatherCondition) error {
	if temp > 128 || temp < -128 {
		return errors.New("temperature out of bounds")
	}
	if temp < 0 {
		temp += 256
	}

	p := getCommandPayload(CmdSetTemperature, []byte{byte(temp), byte(weather)})
	return c.sendMessage(newMessage(p))
}

// SetBrightness controls how bright the Timebox display is
func (c *Conn) SetBrightness(level int) error {
	if level < 0 || level > 100 {
		return errors.New("brightness must be between 0 and 100")
	}

	p := getCommandPayload(CmdSetBrightness, []byte{byte(level)})
	return c.sendMessage(newMessage(p))
}

// SetVolume controls the volume of the Timebox speaker
func (c *Conn) SetVolume(level int) error {
	if level < 0 || level > 100 {
		return errors.New("volume must be between 0 and 100")
	}

	// level is encoded from 0 to 0x10
	level /= 16

	p := getCommandPayload(CmdSetVolume, []byte{byte(level)})
	return c.sendMessage(newMessage(p))
}

// SetMute enables or disables mute on the Timebox speaker
func (c *Conn) SetMute(isMute bool) error {
	var unit byte
	if isMute {
		unit = 0x00
	} else {
		unit = 0x01
	}

	p := getCommandPayload(CmdSetView, []byte{CmdSetMute, unit})
	return c.sendMessage(newMessage(p))

}

// DisplayTemperature sets the Timebox to show the temperature and weather view
func (c *Conn) DisplayTemperature(isCelsius bool) error {
	var unit byte
	if isCelsius {
		unit = 0x00
	} else {
		unit = 0x01
	}

	args := []byte{ViewTemperature, unit}
	if c.colour != nil {
		args = append(args, c.colour.R, c.colour.G, c.colour.B)
	}

	p := getCommandPayload(CmdSetView, args)
	return c.sendMessage(newMessage(p))
}

// DisplayClock sets the Timebox to show the clock
func (c *Conn) DisplayClock(is24Hour bool) error {
	var unit byte
	if is24Hour {
		unit = 0x01
	} else {
		unit = 0x00
	}

	args := []byte{ViewClock, unit}
	if c.colour != nil {
		args = append(args, c.colour.R, c.colour.G, c.colour.B)
	}

	p := getCommandPayload(CmdSetView, args)
	return c.sendMessage(newMessage(p))
}

// DisplayScoreboard shows a scoreboard with the supplied scores
func (c *Conn) DisplayScoreboard(redScore int, blueScore int) error {
	if redScore < 0 {
		redScore = 0
	} else if redScore > 999 {
		redScore = 999
	}
	if blueScore < 0 {
		blueScore = 0
	} else if blueScore > 999 {
		blueScore = 999
	}

	args := []byte{ViewScoreboard, 0}
	args = append(args, encodeInt(redScore)...)
	args = append(args, encodeInt(blueScore)...)

	p := getCommandPayload(CmdSetView, args)
	return c.sendMessage(newMessage(p))
}

// DisplaySolid turns all the lights on to the specified colour
func (c *Conn) DisplaySolid() error {
	args := []byte{ViewSolid}
	if c.colour != nil {
		args = append(args, c.colour.R, c.colour.G, c.colour.B)
	}

	p := getCommandPayload(CmdSetView, args)
	return c.sendMessage(newMessage(p))
}

// DisplayView allows for an arbitrary view with arbitrary args to be set.
// Look to the View* constants in protocol.go, and use the protocol specs to ensure the proper arguments are provided.
func (c *Conn) DisplayView(view byte, args []byte) error {
	p := getCommandPayload(view, args)
	return c.sendMessage(newMessage(p))
}

// TODO: display stopwatch, score board
