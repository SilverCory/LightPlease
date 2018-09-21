package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/huin/goserial"

	"github.com/SilverCory/go-lightpack"
)

var correctionArray = []uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
	2, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 5, 5, 5,
	5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 9, 9, 9, 10,
	10, 10, 11, 11, 11, 12, 12, 13, 13, 13, 14, 14, 15, 15, 16, 16,
	17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22, 22, 23, 24, 24, 25,
	25, 26, 27, 27, 28, 29, 29, 30, 31, 32, 32, 33, 34, 35, 35, 36,
	37, 38, 39, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 50,
	51, 52, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 66, 67, 68,
	69, 70, 72, 73, 74, 75, 77, 78, 79, 81, 82, 83, 85, 86, 87, 89,
	90, 92, 93, 95, 96, 98, 99, 101, 102, 104, 105, 107, 109, 110, 112, 114,
	115, 117, 119, 120, 122, 124, 126, 127, 129, 131, 133, 135, 137, 138, 140, 142,
	144, 146, 148, 150, 152, 154, 156, 158, 160, 162, 164, 167, 169, 171, 173, 175,
	177, 180, 182, 184, 186, 189, 191, 193, 196, 198, 200, 203, 205, 208, 210, 213,
	215, 218, 220, 223, 225, 228, 231, 233, 236, 239, 241, 244, 247, 249, 252, 255}

func main() {

	testing := flag.Bool("test", false, "enables a test mode that cycles through memes.")
	brightest := flag.Bool("brightest", false, "Turn  the light on to it's brightest state.")
	port := flag.String("port", "", "the com/serial port to listen on")
	flag.Parse()

	api := lightpack.API{
		Address: "192.168.0.12:3636",
	}

	if !(*testing) {
		if err := api.Connect(); err != nil {
			panic(err)
			return
		}
	}

	var portName string
	if *port != "" {
		portName = *port
	} else {
		portName = findArduino()
	}

	config := &goserial.Config{Name: portName, Baud: 115200}
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
			sendArduinoCommand("PWM", uint8(0), uint8(0), uint8(0), uint8(0), s)
			s.Close()
			os.Exit(0)
		}
	}()

	fmt.Println("Started and connected.")
	if *testing {
		doTestingMemes(s)
		return
	}

	if *brightest {
		sendArduinoCommand("PWM", uint8(255), uint8(255), uint8(255), uint8(255), s)
		return
	}

	requestNumber := 14
	ledsOn := false
	for {
		// Check if the LEDS are on.
		requestNumber++
		if requestNumber > 15 {
			status, err := api.GetStatus()
			if err != nil {
				fmt.Println(err)
				time.Sleep(3 * time.Second)
			}

			if status != lightpack.StatusOn {
				ledsOn = false
				sendArduinoCommand("PWM", uint8(0), uint8(0), uint8(0), uint8(0), s) // Turn off the LEDs.
			} else {
				ledsOn = true
			}
		}

		if !ledsOn {
			time.Sleep(66 * time.Millisecond) // Sleep 66 millis, because after 15 loops, it will be ~1 second.
			continue
		}

		colors, err := api.GetColors()
		if err != nil {
			fmt.Println(err)
			time.Sleep(3 * time.Second)
		}

		lastColour := colors[len(colors)-2] // Not sure why it's neg 2...
		red, green, blue, white := correctionArray[lastColour.R], correctionArray[lastColour.G], correctionArray[lastColour.B], uint8(0)
		r, g, b, _ := color.RGBA{R: red, G: green, B: blue}.RGBA()
		white = correctionArray[uint8((19595*r+38470*g+7471*b+1<<15)>>24)]

		if err := sendArduinoCommand("PWM", red, green, blue, white, s); err != nil {
			fmt.Println(err)
			if err.Error() != "short write" {
				s, err = goserial.OpenPort(config)
				if err != nil {
					panic(err)
				}
			}
		}

		time.Sleep(23 * time.Millisecond) // Sleep so we don't crash the API.
	}

}

func sendArduinoCommand(command string, red, green, blue, white byte, serialPort io.ReadWriteCloser) error {
	if serialPort == nil {
		return nil
	}

	cmd := fmt.Sprint(command, " ", red, " ", green, " ", blue, " ", white, "\t\n")
	_, err := serialPort.Write([]byte(cmd))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// findArduino looks for the file that represents the Arduino
// serial connection. Returns the fully qualified path to the
// device if we are able to find a likely candidate for an
// Arduino, otherwise an empty string if unable to find
// something that 'looks' like an Arduino device.
func findArduino() string {

	if host, err := os.Hostname(); err != nil && strings.ContainsAny(host, "corys") {
		return "COM10"
	}

	contents, _ := ioutil.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "serial") || strings.Contains(f.Name(), "ttyS") || strings.Contains(f.Name(), "ttyAMA") {
			return "/dev/" + f.Name()
		}
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return ""
}
