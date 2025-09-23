package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"github.com/vasyl-ks/TM-software-H11/config"
	"github.com/vasyl-ks/TM-software-H11/modules/model"
)

func readTelemetry(payload []byte) (model.Telemetry, error) {
    // Decode JSON
    var entry model.Telemetry
    if err := json.Unmarshal(payload, &entry); err != nil {
        return model.Telemetry{}, err
    }
    return entry, nil
}

func logTelemetry(entry model.Telemetry) error {
    fileDir := config.Client.FileDir
    // Check directory
    err := os.MkdirAll(fileDir, 0755)
    if err != nil {
        fmt.Println("Error in creating directory:", err)
        return err
    }

    // Create a new file
    filename := fmt.Sprintf("%s/telemetry_%s.log", fileDir, time.Now().Format("20060102_15405.000"))
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write in the file
    log.SetOutput(file)
    log.SetFlags(0)
    log.Printf("| Created at %s, Logged at %s | Vehicle:%s, Speed:%.2f, RPM=%.2f, Temp=%.2f, Pressure=%.2f\n",
        entry.CreatedAt.Local().Format("15:04:05.000000"),
        time.Now().Local().Format("15:04:05.000000"),
        entry.VehicleID,
        entry.Speed,
        entry.RPM,
        entry.Temperature,
        entry.Pressure,
    )
    return nil
}

func processDatagrams(addr *net.UDPAddr, payload []byte) {
    // Parse datagram
	entry, err := readTelemetry(payload) 
	if err != nil {
		fmt.Printf("Datagram decoding from %s failed: %v\n", addr.String(), err)
		return
	}

    // Log datagram
	if err := logTelemetry(entry); err != nil {
		fmt.Printf("Datagram log for %s failed: %v\n", addr.String(), err)
		return
	}
}

func UDPClient() {
	addr := net.UDPAddr{Port: config.Server.ClientPort} // UDP Address
	
    // Listen for UDP Traffic
    conn, err := net.ListenUDP("udp", &addr) 
    if err != nil {
		fmt.Println("Error listening:", err)
    }
    defer conn.Close()
	fmt.Printf("Listening on %s", conn.LocalAddr())
	
    buf := make([]byte, 1024) // Buffer for incoming datagrams
    for {
        // Read one datagram and track the sender
        n, remoteAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Printf("read UDP: %v", err)
            continue
        }

        // Clone current datagram
        payload := make([]byte, n)
        copy(payload, buf[:n])

        // Handle this datagram concurrently (Parse and Log)
        go processDatagrams(remoteAddr, payload)
    }
}