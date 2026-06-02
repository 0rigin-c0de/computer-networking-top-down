package main

import (
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	go startSniffer()

	time.Sleep(1 * time.Second)

	listener, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Printf("Listen Error: %v\n", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go func(c net.Conn) {
			time.Sleep(200 * time.Millisecond)
			c.Close()
		}(conn)
	}
}

func startSniffer() {
	handle, err := pcap.OpenLive("lo0", 1024, false, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Sniffer Interface Error: %v\n", err)
		return
	}
	defer handle.Close()

	err = handle.SetBPFFilter("tcp port 9000")
	if err != nil {
		fmt.Printf("BPF Filter Error: %v\n", err)
		return
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	fmt.Println("--- Sniffer Active: Tracking Port 9000 State Machine ---")

	for packet := range packetSource.Packets() {
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}

		tcp, _ := tcpLayer.(*layers.TCP)

		var flags []string
		if tcp.SYN {
			flags = append(flags, "SYN")
		}
		if tcp.ACK {
			flags = append(flags, "ACK")
		}
		if tcp.FIN {
			flags = append(flags, "FIN")
		}
		if tcp.RST {
			flags = append(flags, "RST")
		}

		if len(flags) > 0 {
			fmt.Printf("[TCP Packet] Src: %d -> Dst: %d | Flags: %v | Seq: %d | Ack: %d\n",
				tcp.SrcPort, tcp.DstPort, flags, tcp.Seq, tcp.Ack)
		}
	}
}
