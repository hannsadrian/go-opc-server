package main

import (
	"./dao"
	"./domain"
	"fmt"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"net"
	"os"
)

func main() {
	// TODO: refactor the Ws2811Driver loading to be based on settings provided by a config file!
	opt := domain.Ws2811DefaultOptions
	engine, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		panic(err)
	}
	driver := domain.Ws2811Driver{
		Options: opt,
		Engine:  engine,
	}


	err = driver.Setup()
	if err != nil {
		panic(err)
	}
	defer driver.Unregister()


	// listen on all interfaces
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
		go handleOPCRequest(&driver, conn)
	}
}

func handleOPCRequest(driver dao.Driver, conn net.Conn) {
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
		go driver.Display(channel, pixels)
	}
}
