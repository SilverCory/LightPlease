package main

import (
	"fmt"
	"math"
	"time"

	"github.com/SilverCory/LightPlease"
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

	io := LightPlease.NewIOOut(12, 13, 18, 19)

	for {
		colors, err := api.GetColors()
		if err != nil {
			fmt.Println(err)
			time.Sleep(3 * time.Second)
		}

		var R, G, B, W int16
		for _, color := range colors {
			R = R + (color.R ^ 2)
			G = G + (color.G ^ 2)
			B = B + (color.B ^ 2)
		}
		W = 0

		io.DisplayRGBW(
			int16(math.Sqrt(float64(R))),
			int16(math.Sqrt(float64(G))),
			int16(math.Sqrt(float64(B))),
			int16(math.Sqrt(float64(W))),
		)
	}

}
