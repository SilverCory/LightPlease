package main

import (
	"fmt"
	"io"
	"time"
)

var HSVlights = [61]uint8{0, 4, 8, 13, 17, 21, 25, 30, 34, 38, 42, 47, 51, 55, 59, 64, 68, 72, 76,
	81, 85, 89, 93, 98, 102, 106, 110, 115, 119, 123, 127, 132, 136, 140, 144,
	149, 153, 157, 161, 166, 170, 174, 178, 183, 187, 191, 195, 200, 204, 208,
	212, 217, 221, 225, 229, 234, 238, 242, 246, 251, 255}

func doTestingMemes(s io.ReadWriteCloser) {

	for {

		fadeInAndOut(func(r int) {
			if err := sendArduinoCommand("PWM", correctionArray[r], 0, 0, 0, s); err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Millisecond * 15)
		})

		fadeInAndOut(func(g int) {
			if err := sendArduinoCommand("PWM", 0, correctionArray[g], 0, 0, s); err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Millisecond * 15)
		})

		fadeInAndOut(func(b int) {
			if err := sendArduinoCommand("PWM", 0, 0, correctionArray[b], 0, s); err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Millisecond * 15)
		})

		fadeInAndOut(func(w int) {
			if err := sendArduinoCommand("PWM", 0, 0, 0, byte(w), s); err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Millisecond * 15)
		})

		for i := 0; i < 15; i++ {
			for w := 0; w < 360; w++ {
				r, g, b := createColourMatrix(w)
				if err := sendArduinoCommand("PWM", r, g, b, 0, s); err != nil {
					fmt.Println(err)
				}
				time.Sleep(10 * time.Millisecond)
			}
		}

	}

}

func fadeInAndOut(callBack func(out int)) {
	for i := 0; i < len(correctionArray); i++ {
		callBack(i)
	}

	for i := len(correctionArray) - 1; i > 0; i-- {
		callBack(i)
	}
}

func createColourMatrix(angle int) (red, green, blue byte) {
	angle = angle + 15
	if angle > 360 {
		angle = angle - 360
	}

	if angle < 60 {
		red = 255
		green = HSVlights[angle]
		blue = 0
	} else if angle < 120 {
		red = HSVlights[120-angle]
		green = 255
		blue = 0
	} else if angle < 180 {
		red = 0
		green = 255
		blue = HSVlights[angle-120]
	} else if angle < 240 {
		red = 0
		green = HSVlights[240-angle]
		blue = 255
	} else if angle < 300 {
		red = HSVlights[angle-240]
		green = 0
		blue = 255
	} else {
		index := 360 - angle
		if index < 0 || index > len(HSVlights) {
			index = 0
		}
		red = 255
		green = 0
		blue = HSVlights[index]
	}

	return
}
