package main

import (
	"fmt"
	"net"
	"os"
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

var DefaultOptions = ws2811.Option{
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

func main() {
	opt := DefaultOptions

	engine, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		panic(err)
	}

	err = engine.Init()
	if err != nil {
		panic(err)
	}
	defer engine.Fini()


	//listen on all interfaces
	l, err := net.Listen("tcp", "0.0.0.0:7890")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + "0.0.0.0" + ":" + "7890")
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleOPCRequest(engine, conn)
	}
}

func handleOPCRequest(engine wsEngine, conn net.Conn) {
	for {
		buf := make([]byte, 65539)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Connection closing:", err.Error())
			return
		}
		channel := int(buf[0])
		command := int(buf[1]) // FIXME: currently only command 0 for 8-bit pixel colors is supported!
		length := (int(buf[2]) << 8) + int(buf[3])
		pixels := make([][]int, 0, 21845)
		fmt.Printf("Chan: %d, Cmd: %d, Len: %d \n", channel, command, length)
		for i := 0; i < length; i += 3 {
			r := int(buf[4+i+0])
			g := int(buf[4+i+1])
			b := int(buf[4+i+2])
			pixels = append(pixels, []int{r, g, b})
		}
		go displayPixels(engine, channel, pixels)
	}
}

func displayPixels(engine wsEngine, channel int, pixels [][]int) {
	for i := 0; i < len(engine.Leds(channel)) && i < len(pixels); i++ {
		engine.Leds(channel)[i] = rgbToColor(pixels[i])
		_ = engine.Render()
	}
}

func rgbToColor(c []int) uint32 {
	i, err := strconv.ParseUint(fmt.Sprintf("%02x", c[0]) + fmt.Sprintf("%02x", c[1]) + fmt.Sprintf("%02x", c[2]), 16, 32)
	if err != nil {
		return 0
	}
	return uint32(i)
}