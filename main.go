package main

import (
	"fmt"
	"github.com/adwirawien/go-opc-server/dao"
	"github.com/adwirawien/go-opc-server/domain"
	"net"
	"os"
)

func main() {
	// setup domain object
	driver := domain.Ws281xDriver{}
	err := driver.Setup(domain.Ws281xDefaultOptions)
	if err != nil {
		panic(err)
	}
	defer driver.Unregister()

	// listen on all interfaces
	address := "0.0.0.0:7890"
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Listening on %s", address)
	defer l.Close()

	for {
		// accept new connections and parse them
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleOPCRequest(&driver, conn)
	}
}

// handleOPCRequest is responsible for handling an incoming connection and the transmitted data respectively
func handleOPCRequest(driver dao.Driver, conn net.Conn) {
	for {
		buf := make([]byte, 65539)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Connection closing:", err.Error())
			return
		}

		// parse header bytes with opc information
		channel := int(buf[0])
		command := int(buf[1]) // FIXME: currently only command 0 for 8-bit pixel colors is supported!
		length := (int(buf[2]) << 8) + int(buf[3])

		// parse pixel information
		pixels := make([][]int, 0, 21845)
		fmt.Printf("Chan: %d, Cmd: %d, Len: %d \n", channel, command, length)
		for i := 0; i < length; i += 3 {
			// use 4 byte offset due to the header bytes at the buffer start
			r := int(buf[4+i+0])
			g := int(buf[4+i+1])
			b := int(buf[4+i+2])
			pixels = append(pixels, []int{r, g, b})
		}

		go driver.Display(channel, pixels)
	}
}
