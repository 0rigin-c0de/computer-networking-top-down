package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	dnsServer := "8.8.8.8:53"

	conn, err := net.Dial("udp", dnsServer)
	if err != nil {
		fmt.Printf("couldnt connect: %v\n", err)
		return
	}
	defer conn.Close()

	packet := make([]byte, 12)

	packet[0], packet[1] = 0xAA, 0xBB
	packet[2] = 0x01
	packet[3] = 0x00
	packet[4], packet[5] = 0x00, 0x01
	packet[6], packet[7] = 0x00, 0x00
	packet[8], packet[9] = 0x00, 0x00
	packet[10], packet[11] = 0x00, 0x00

	// domain name has to be split by dots and each part needs length prefix i think
	domain := "example.com"
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		packet = append(packet, byte(len(part)))
		packet = append(packet, []byte(part)...)
	}
	packet = append(packet, 0x00)

	// type A and class IN
	packet = append(packet, 0x00, 0x01)
	packet = append(packet, 0x00, 0x01)

	_, err = conn.Write(packet)
	if err != nil {
		fmt.Printf("write failed: %v\n", err)
		return
	}

	reply := make([]byte, 512)
	n, err := conn.Read(reply)
	if err != nil {
		fmt.Printf("read failed: %v\n", err)
		return
	}

	fmt.Printf("sent %d bytes\n", len(packet))
	fmt.Printf("got back %d bytes\n", n)
	fmt.Println("--------------------------------------------------")

	fmt.Printf("transaction id from response: %02X%02X\n", reply[0], reply[1])

	// last 4 bytes should be the ip hopefully
	if n >= 4 {
		ip := reply[n-4 : n]
		fmt.Printf("ip for %s is: %d.%d.%d.%d\n", domain, ip[0], ip[1], ip[2], ip[3])
	}
}
