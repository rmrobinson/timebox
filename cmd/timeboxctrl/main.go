package main

import (
	"flag"
	"log"
	"time"

	"github.com/rmrobinson/timebox"
	bt "github.com/rmrobinson/timebox/bluetooth"
)

const (
	timeboxMiniBluetoothChannel = 4
)

var (
	timeboxAddr = flag.String("addr", "", "The Bluetooth address to connect to")
	view        = flag.String("view", "", "The view to show. Clock, Temperature, Scoreboard, Solid")
	arg         = flag.Int("arg", 0, "arg")
)

func main() {
	flag.Parse()

	btAddr, err := bt.NewAddress(*timeboxAddr)
	if err != nil {
		log.Fatalf("invalid bluetooth address (%s): %s\n", *timeboxAddr, err)
	}

	btConn := &bt.Connection{}
	err = btConn.Connect(btAddr, timeboxMiniBluetoothChannel)
	if err != nil {
		log.Fatalf("unable to connect to bluetooth device: %s\n", err.Error())
	}
	defer btConn.Close()

	tbConn := timebox.NewConn(btConn)
	if err := tbConn.Initialize(); err != nil {
		log.Fatalf("unable to establish connection with timebox: %s\n", err.Error())
	}

	tbConn.SetColor(&timebox.Colour{R: 0, G: 255, B: 66})

	tbConn.SetBrightness(100)
	tbConn.SetTime(time.Now())
	tbConn.SetTemperatureAndWeather(37, timebox.Celsius, timebox.WeatherCondition(*arg))

	if *view == "Clock" {
		tbConn.DisplayClock(true)
	} else if *view == "Temperature" {
		tbConn.DisplayTemperature(true)
	} else if *view == "Scoreboard" {
		tbConn.DisplayScoreboard(10, 20)
	} else if *view == "Solid" {
		tbConn.DisplaySolid()
	} else if *view == "Off" {
		tbConn.SetBrightness(0)
	} else {
		log.Printf("view not supported")
	}
}
