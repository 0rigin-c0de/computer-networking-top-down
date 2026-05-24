package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
)

func main() {
	port := ":9090"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		fmt.Println("something went wrong:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("couldn't start server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("server up on", port)

	buf := make([]byte, 1024)
	for {
		n, from, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		msg := string(buf[:n])
		fmt.Printf("got %q from %s\n", msg, from)

		if rand.Intn(10) < 3 {
			fmt.Println("dropped")
			continue
		}

		_, err = conn.WriteToUDP([]byte(msg), from)
		if err != nil {
			fmt.Println("send error:", err)
		}
	}
}
