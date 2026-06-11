package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	target := "8.8.8.8"
	dstAddr, err := net.ResolveIPAddr("ip4", target)
	if err != nil {
		fmt.Printf("Resolution error: %v\n", err)
		return
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Printf("Socket error (Run with sudo): %v\n", err)
		return
	}
	defer conn.Close()

	packetConn := conn.IPv4PacketConn()
	maxHops := 30
	timeout := 2 * time.Second

	fmt.Printf("Tracing route to %s [%s] over a maximum of %d hops:\n\n", target, dstAddr.IP, maxHops)

	for ttl := 1; ttl <= maxHops; ttl++ {
		// Explicitly restrict the IP layer package lifetime
		if err := packetConn.SetTTL(ttl); err != nil {
			fmt.Printf("Error setting TTL: %v\n", err)
			return
		}

		// Construct standard Echo Request payload
		msg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  ttl,
				Data: []byte("NET-TRACE"),
			},
		}

		msgBytes, err := msg.Marshal(nil)
		if err != nil {
			return
		}

		startTime := time.Now()
		if _, err := conn.WriteTo(msgBytes, dstAddr); err != nil {
			fmt.Printf("%d\tError sending packet\n", ttl)
			continue
		}

		_ = conn.SetReadDeadline(time.Now().Add(timeout))

		replyBuffer := make([]byte, 1500)
		n, peer, err := conn.ReadFrom(replyBuffer)
		duration := time.Since(startTime)

		if err != nil {
			// Router configuration drop or silent drop timeout
			fmt.Printf("%d\t*\tRequest timed out.\n", ttl)
			continue
		}

		parsedMsg, err := icmp.ParseMessage(1, replyBuffer[:n])
		if err != nil {
			fmt.Printf("%d\tError parsing packet from %s\n", ttl, peer.String())
			continue
		}

		switch parsedMsg.Type {
		case ipv4.ICMPTypeTimeExceeded:
			// Packet layer hit 0 TTL at intermediate hop
			fmt.Printf("%d\t%d ms\t%s\n", ttl, duration.Milliseconds(), peer.String())
		case ipv4.ICMPTypeEchoReply:
			// Reached final destination target
			fmt.Printf("%d\t%d ms\t%s (Reached Target)\n", ttl, duration.Milliseconds(), peer.String())
			return
		default:
			fmt.Printf("%d\tUncaught Type: %v from %s\n", ttl, parsedMsg.Type, peer.String())
		}
	}
}
