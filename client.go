package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var addr = flag.String("server", "198.100.31.182:9999", "server address")
var duration = flag.Int("duration", -1, "idle duration in minutes")
var defaultDuration = 15

func askDuration() int {
	var d int
	fmt.Printf("Give idle time in minutes or press enter to use default [%d]: ", defaultDuration)
	_, err := fmt.Scanln(&d)
	if err != nil {
		return defaultDuration
	}
	return d
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *duration == -1 {
		*duration = askDuration()
	}

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalln("Could not connect to server:", err)
	}

	log.Println("Connection established.")

	response := bufio.NewReader(conn)

	log.Printf("Now idling %d minutes before testing connection, please wait...\n", *duration)

	time.Sleep(time.Duration(*duration) * time.Minute)

	// idle for extra 10 seconds just to make sure
	time.Sleep(10 * time.Second)

	// set a deadline so we don't have to wait too much
	deadline := 3 * time.Second
	conn.SetDeadline(time.Now().Add(deadline))

	msg := []byte("A")
	_, err = conn.Write(msg)
	if err != nil {
		log.Fatalln("Connection error:", err)
	}

	for {
		b, err := response.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Connection error:", err)
		}
		if b == 65 {
			log.Printf("Connection works fine!")
		}

	}
	conn.Close()
}
