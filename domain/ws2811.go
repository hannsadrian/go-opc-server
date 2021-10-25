package domain

import (
	"fmt"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"strconv"
)

var gamma8 = []byte{
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
	215, 218, 220, 223, 225, 228, 231, 233, 236, 239, 241, 244, 247, 249, 252, 255,
}

type wsEngine interface {
	Init() error
	Render() error
	Wait() error
	Fini()
	Leds(channel int) []uint32
}

var Ws2811DefaultOptions = ws2811.Option{
	Frequency: ws2811.TargetFreq,
	DmaNum:    ws2811.DefaultDmaNum,
	Channels: []ws2811.ChannelOption{
		{
			GpioPin:    18,
			LedCount:   10,
			Brightness: 255,
			StripeType: ws2811.WS2812Strip,
			Invert:     false,
			Gamma:      gamma8,
		},
		{
			GpioPin:    13,
			LedCount:   10,
			Brightness: 255,
			StripeType: ws2811.WS2812Strip,
			Invert:     false,
			Gamma:      gamma8,
		},
	},
}

type Ws2811Driver struct {
	Options ws2811.Option
	Engine wsEngine
}

func (d *Ws2811Driver) GetInstance() interface{} {
	return d
}

func (d *Ws2811Driver) Setup() error {
	err := d.Engine.Init()
	if err != nil {
		return err
	}
	return nil
}

func (d *Ws2811Driver) Display(channel int, pixels [][]int) {
	displayPixels(d.Engine, channel, pixels)
}

func (d *Ws2811Driver) Unregister() {
	d.Engine.Fini()
}

func displayPixels(engine wsEngine, channel int, pixels [][]int) {
	for i := 0; i < len(engine.Leds(channel)) && i < len(pixels); i++ {
		engine.Leds(channel)[i] = rgbToColor(pixels[i])
		_ = engine.Render()
	}
}

func rgbToColor(c []int) uint32 {
	i, err := strconv.ParseUint(
		fmt.Sprintf("%02x", c[0]) +
		fmt.Sprintf("%02x", c[1]) +
		fmt.Sprintf("%02x", c[2]),
		16,
		32,
	)
	if err != nil {
		return 0
	}
	return uint32(i)
}
