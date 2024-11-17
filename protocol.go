package timebox

const (
	CmdSwitchRadio       = 0x05
	CmdSetVolume         = 0x08
	CmdGetVolume         = 0x09
	CmdSetMute           = 0x0a
	CmdGetMute           = 0x0b
	CmdSetDateTime       = 0x18
	CmdSetView           = 0x45
	CmdSetAnimationFrame = 0x49
	CmdGetTemperature    = 0x59
	CmdSetTemperature    = 0x5f
	CmdGetRadioFrequency = 0x60
	CmdSetRadioFrequency = 0x61
	CmdSetBrightness     = 0x74
)

const (
	ViewClock       = 0x00
	ViewTemperature = 0x01
	ViewSolid       = 0x02
	ViewAnimation   = 0x03
	ViewGraph       = 0x04
	ViewImage       = 0x05
	ViewStopwatch   = 0x06
	ViewScoreboard  = 0x07
)

const (
	WeatherSun                  = 1
	WeatherSunAndCloud          = 2
	WeatherPartiallyCloudy      = 3
	WeatherCloudy               = 4
	WeatherRain                 = 5
	WeatherRainAndSun           = 6
	WeatherRainAndLightning     = 7
	WeatherSnow                 = 8
	WeatherFog                  = 9
	WeatherDarkClear            = 10
	WeatherDarkPartiallyCoudy   = 11
	WeatherDarkCloudy           = 12
	WeatherDarkRain             = 13
	WeatherDarkRainAndSun       = 14
	WeatherDarkRainAndLightning = 15
	WeatherDarkSnow             = 16
	WeatherDarkFog              = 17
)

var helloResponse = []byte{0, 5, 72, 69, 76, 76, 79, 0}

// getCommandPayload converts a specified command and its arguments into a single payload for further encoding ahead of transmission
func getCommandPayload(cmd byte, args []byte) []byte {
	payloadLen := len(args) + 3

	lenLsb := byte(payloadLen & 0xff)
	lenMsb := byte(payloadLen >> 8)

	payload := []byte{lenLsb, lenMsb, cmd}
	if len(args) > 0 {
		payload = append(payload, args...)
	}

	return payload
}

// checksum calculated of the payload. Returned as the LSB, MSB for use in subsequent functions
func checksum(payload []byte) []byte {
	csum := 0

	for _, b := range payload {
		csum += int(b)
	}

	lsb := byte(csum & 0x00ff)
	msb := byte(csum >> 8)

	return []byte{lsb, msb}
}

// escape the special characters (0x01, 0x02 and 0x03) in the supplied byte array.
func escape(payload []byte) []byte {
	var escaped []byte

	for _, b := range payload {
		if b == 0x01 || b == 0x02 || b == 0x03 {
			escaped = append(escaped, 0x03)
			escaped = append(escaped, byte(b+0x03))
		} else {
			escaped = append(escaped, byte(b))
		}
	}

	return escaped
}

// unescape the special characters (0x01, 0x02 and 0x03) in the supplied byte array.
func unescape(payload []byte) ([]byte, error) {
	var ret []byte
	wasEscaped := false

	for b := range payload {
		if wasEscaped {
			if b < 0x04 || b > 0x06 {
				return nil, ErrMalformedPayload
			}
			ret = append(ret, byte(b-0x03))
			wasEscaped = false
			continue
		}

		if b == 0x03 {
			wasEscaped = true
		} else {
			ret = append(ret, byte(b))
		}
	}

	return ret, nil
}

func encodeInt(val int) []byte {
	lsb := byte(val & 0x00ff)
	msb := byte(val >> 8)

	return []byte{lsb, msb}
}
