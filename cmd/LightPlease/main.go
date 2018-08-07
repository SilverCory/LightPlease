package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/huin/goserial"

	"github.com/SilverCory/go-lightpack"
)

func main() {
	api := lightpack.API{
		Address: "192.168.0.12:3636",
	}

	if err := api.Connect(); err != nil {
		panic(err)
		return
	}

	config := &goserial.Config{Name: findArduino(), Baud: 115200}
	s, err := goserial.OpenPort(config)
	if err != nil {
		panic(err)
	}

	// When connecting to an older revision Arduino, you need to wait
	// a little while it resets.
	time.Sleep(1 * time.Second)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			s.Close()
			os.Exit(0)
		}
	}()

	for {
		colors, err := api.GetColors()
		if err != nil {
			fmt.Println(err)
			time.Sleep(3 * time.Second)
		}

		var R, G, B, W uint16
		for _, color := range colors {
			R = R + (uint16(color.R) ^ 2)
			G = G + (uint16(color.G) ^ 2)
			B = B + (uint16(color.B) ^ 2)
		}
		W = 0 + W

		if err := sendArduinoCommand(byte('F'), AverageColourValue(R, len(colors)), AverageColourValue(G, len(colors)), AverageColourValue(B, len(colors)), AverageColourValue(W, len(colors)), s); err != nil {
			fmt.Println(err)
			s, err = goserial.OpenPort(config)
			if err != nil {
				panic(err)
			}
		}
	}

}

func AverageColourValue(colour uint16, length int) uint8 {
	return uint8(math.Floor(math.Sqrt(float64(int(colour) / length))))
}

func sendArduinoCommand(command byte, red, green, blue, white byte, serialPort io.ReadWriteCloser) error {
	if serialPort == nil {
		return nil
	}

	// Transmit command and argument down the pipe.
	for _, v := range [][]byte{{command}, {red}, {green}, {blue}, {white}} {
		_, err := serialPort.Write(v)
		if err != nil {
			return err
		}
	}

	return nil
}

// findArduino looks for the file that represents the Arduino
// serial connection. Returns the fully qualified path to the
// device if we are able to find a likely candidate for an
// Arduino, otherwise an empty string if unable to find
// something that 'looks' like an Arduino device.
func findArduino() string {
	contents, _ := ioutil.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "ttyUSB") {
			return "/dev/" + f.Name()
		}
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return ""
}
