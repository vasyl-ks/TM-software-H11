package consumer

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/vasyl-ks/TM-software-H11/config"
)

var Ready = make(chan struct{})

/*
Listen binds a UDP socket on config.Sender.ClientPort and forwards incoming datagrams to out.
- Copies each datagram into a new slice to avoid buffer reuse.
*/
func Listen(outChan chan<- []byte) {
	addrUDP := fmt.Sprintf("127.0.0.1:%d", config.Hub.UDPPort)
	addrTCP := fmt.Sprintf("127.0.0.1:%d", config.Hub.TCPPort)

	// Resolve and bind UDP
	udpAddr, err := net.ResolveUDPAddr("udp", addrUDP)
	if err != nil {
		log.Printf("[ERROR][Consumer][Listen] Failed to resolve UDP address %s: %v", addrUDP, err)
		return
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("[ERROR][Consumer][Listen] Failed to listen on UDP %s: %v", addrUDP, err)
		return
	}

	// Bind TCP
	tcpListener, err := net.Listen("tcp", addrTCP)
	if err != nil {
		log.Printf("[ERROR][Consumer][Listen] Failed to listen on TCP %s: %v", addrTCP, err)
		return
	}

	// Notify that listeners are ready
	close(Ready)
	log.Printf("[INFO][Consumer][Listen] Listening on UDP %s and TCP %s", addrUDP, addrTCP)

	// UDP handler goroutine
	go func() {
		defer udpConn.Close()
		buf := make([]byte, config.Hub.BufferSize)
		for {
			n, _, err := udpConn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("[ERROR][Consumer][Listen] Error reading from UDP: %v\n", err)
				continue
			}
			payload := make([]byte, n)
			copy(payload, buf[:n])
			outChan <- payload
		}
	}()

	// TCP handler goroutine
	go func ()  {
		defer tcpListener.Close()
		conn, err := tcpListener.Accept()
			if err != nil {
			fmt.Println("[ERROR][Consumer][Listen] Error accepting TCP connection:", err)
			return
		}
		defer conn.Close()
		buf := make([]byte, config.Hub.BufferSize)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("[ERROR][Consumer][Listen] Error reading TCP: %v\n", err)
				}
				break
			}
			payload := make([]byte, n)
			copy(payload, buf[:n])
			outChan <- payload
		}
	}()
}
