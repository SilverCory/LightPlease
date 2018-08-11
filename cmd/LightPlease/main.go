package main

import (
	"fmt"
	col "image/color"
	"io"
	"io/ioutil"
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

		var R, G, B, W uint32
		for _, color := range colors {
			rgbaCol := col.RGBA{R: color.R, G: color.G, B: color.B, A: 255}
			r, g, b, _ := rgbaCol.RGBA()
			R += r
			G += g
			B += b
		}

		length := uint32(len(colors))
		R /= length
		G /= length
		B /= length

		if err := sendArduinoCommand(byte('F'), uint8(R), uint8(G), uint8(B), uint8(W), s); err != nil {
			fmt.Println(err)
			if err.Error() != "short write" {
				s, err = goserial.OpenPort(config)
				if err != nil {
					panic(err)
				}
			}
		}
		time.Sleep(5 * time.Millisecond)
	}

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
			strings.Contains(f.Name(), "ttyUSB") || strings.Contains(f.Name(), "serial") {
			return "/dev/" + f.Name()
		}
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return ""
}
