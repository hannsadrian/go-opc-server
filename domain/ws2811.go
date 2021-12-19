package domain

import (
	"fmt"
	ws281x "github.com/rpi-ws281x/rpi-ws281x-go"
)

// gamma8 is a calibration byte array for the LEDs
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

// wsEngine represents the capabilities of the ws281x library
type wsEngine interface {
	Init() error
	Render() error
	Wait() error
	Fini()
	Leds(channel int) []uint32
}

// Ws281xDefaultOptions holds configuration for the ws281x library
// that is adjusted to a setup of two led strands with a length of ten pixels
var Ws281xDefaultOptions = ws281x.Option{
	Frequency: ws281x.TargetFreq,
	DmaNum:    ws281x.DefaultDmaNum,
	Channels: []ws281x.ChannelOption{
		{
			GpioPin:    18,
			LedCount:   10,
			Brightness: 255,
			StripeType: ws281x.WS2812Strip,
			Invert:     false,
			Gamma:      gamma8,
		},
		{
			GpioPin:    13,
			LedCount:   10,
			Brightness: 255,
			StripeType: ws281x.WS2812Strip,
			Invert:     false,
			Gamma:      gamma8,
		},
	},
}

// Ws281xDriver implements the dao.Driver interface along with
// custom properties specific to the ws281x library
type Ws281xDriver struct {
	Options ws281x.Option
	Engine wsEngine
}

// GetInstance returns the ws281xDriver instance
func (d *Ws281xDriver) GetInstance() interface{} {
	return d
}

// Setup initializes the ws281x library
func (d *Ws281xDriver) Setup(untypedOptions interface{}) error {
	options, ok := untypedOptions.(ws281x.Option)
	if !ok {
		return fmt.Errorf("given interface could not be parsed to ws281x.Option")
	}

	engine, err := ws281x.MakeWS2811(&options)
	if err != nil {
		return err
	}
	d.Options = options
	d.Engine = engine

	err = d.Engine.Init()
	if err != nil {
		return err
	}
	return nil
}

// Display calls displayPixels to change the state of the LED's
func (d *Ws281xDriver) Display(channel int, pixels [][]int) {
	displayPixels(d.Engine, channel, pixels)
}

// Unregister takes care of cleanup for the ws281x library
func (d *Ws281xDriver) Unregister() {
	d.Engine.Fini()
}

// displayPixels interfaces with wsEngine to display pixel data
func displayPixels(engine wsEngine, channel int, pixels [][]int) {
	for i := 0; i < len(engine.Leds(channel)) && i < len(pixels); i++ {
		engine.Leds(channel)[i] = rgbToHex(pixels[i])
		_ = engine.Render()
	}
}

// rgbToHex converts an array of [r,g,b] values from 0 to 255
// into a hex value represented by an uint32
func rgbToHex(rgb []int) uint32 {
    return uint32(rgb[0])<<16 | uint32(rgb[1])<<8 | uint32(rgb[2])
}
