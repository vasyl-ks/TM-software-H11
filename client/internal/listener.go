package internal

import (
	"fmt"
	"net"

	"github.com/vasyl-ks/TM-software-H11/config"
)

/*
Listen binds a UDP socket on config.Sender.ClientPort and forwards incoming datagrams to out.
- Copies each datagram into a new slice to avoid buffer reuse.
-Intended to run as a goroutine and loops indefinitely until shutdown.
*/
func Listen(out chan<- []byte) {
	addr := net.UDPAddr{Port: config.Sender.ClientPort} // UDP Address

	// Listen for UDP Traffic
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
	}
	defer conn.Close()
	fmt.Printf("Listening on %s", conn.LocalAddr())

	buf := make([]byte, 1024) // Buffer for incoming datagram
	for {
		// Read one datagram and track the sender
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("read UDP: %v", err)
			continue
		}

		// Clone current datagram
		payload := make([]byte, n)
		copy(payload, buf[:n])

		// Send the datagram through the channel
		out <- payload
	}
}
