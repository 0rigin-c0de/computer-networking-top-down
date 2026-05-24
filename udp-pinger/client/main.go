package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	server := "127.0.0.1:9090"
	if len(os.Args) > 1 {
		server = os.Args[1]
	}

	addr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		fmt.Println("bad address:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("couldn't connect:", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	lost := 0
	var total time.Duration

	for i := 1; i <= 10; i++ {
		msg := fmt.Sprintf("ping %d", i)
		start := time.Now()

		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("send failed:", err)
			os.Exit(1)
		}

		conn.SetReadDeadline(time.Now().Add(time.Second))

		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("request %d timed out\n", i)
			lost++
			time.Sleep(time.Second)
			continue
		}

		rtt := time.Since(start)
		total += rtt
		fmt.Printf("reply from %s: %s time=%.2fms\n", addr, string(buf[:n]), float64(rtt.Microseconds())/1000)
		time.Sleep(time.Second)
	}

	received := 10 - lost
	fmt.Printf("\npackets: sent=10 received=%d lost=%d loss=%.0f%%\n", received, lost, float64(lost)*10)
	if received > 0 {
		avg := total / time.Duration(received)
		fmt.Printf("average time: %.2fms\n", float64(avg.Microseconds())/1000)
	}
}
