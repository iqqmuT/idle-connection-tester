package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
)

var addr = flag.String("addr", "localhost:9999", "service address")

func echo(conn net.Conn) {
	defer conn.Close()
	// read & write twice per connection
	for i := 0; i < 2; i++ {
		r := bufio.NewReader(conn)
		b, err := r.ReadByte()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Println("ERROR", err)
		}
		log.Println("Received:", b)
		resp := make([]byte, 1)
		resp[0] = b
		conn.Write(resp)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalln("Listen error - ", err)
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
		go echo(conn)
	}
}
