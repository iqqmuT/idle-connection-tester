package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

func quit(code int) {
	// wait for enter keypress
	log.Println("Press enter to quit.")
	var i int
	fmt.Scanln(&i)
	os.Exit(code)
}

func printSuccess() {
	log.Println("")
	log.Println("  +------------------------+")
	log.Println("  | CONNECTION WORKS FINE! |")
	log.Println("  +------------------------+")
	log.Println("")
}

func printFailure() {
	log.Println("")
	log.Println("  +--------------------+")
	log.Println("  | CONNECTION FAILED! |")
	log.Println("  +--------------------+")
	log.Println("")
}

func pretest(conn net.Conn, reader *bufio.Reader) {
	msg := []byte("A")
	_, err := conn.Write(msg)
	if err != nil {
		log.Println("Error when pre-testing connection:", err)
		conn.Close()
		quit(1)
	}

	b, err := reader.ReadByte()
	if err == io.EOF {
		log.Println("Error when pre-testing connection: server closed connection")
		conn.Close()
		quit(1)
	}
	if err != nil {
		log.Println("Error when pre-testing connection:", err)
		conn.Close()
		quit(1)
	}
	if b == 65 {
		log.Println("Connection established and tested.")
	}
}

func test(conn net.Conn, reader *bufio.Reader) {
	log.Printf("Idling %d minute(s) before testing connection again, please wait...\n", *duration)

	time.Sleep(time.Duration(*duration) * time.Minute)

	// idle for extra 10 seconds just to make sure
	time.Sleep(10 * time.Second)

	log.Println("Sending message and waiting for response...")

	// set a deadline so we don't have to wait for minutes
	deadline := 15 * time.Second
	conn.SetDeadline(time.Now().Add(deadline))

	msg := []byte("B")
	_, err := conn.Write(msg)
	if err != nil {
		printFailure()
		log.Println("Connection error:", err)
		quit(1)
	}

	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			printFailure()
			log.Println("Connection error:", err)
			quit(1)
		}
		if b == 66 {
			printSuccess()
		} else {
			log.Println("Response was", b)
		}

	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *duration == -1 {
		*duration = askDuration()
	}

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Println("Could not connect to server:", err)
		quit(1)
	}

	reader := bufio.NewReader(conn)

	pretest(conn, reader)
	test(conn, reader)

	conn.Close()
	quit(0)
}
