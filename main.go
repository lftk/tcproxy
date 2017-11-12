package main

import (
	"flag"
	"log"
	"net"
)

var (
	addr    = flag.String("addr", ":9090", "tcproxy address")
	forward = flag.String("forward", "", "forwarding address")
)

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go func(c net.Conn) {
			to, err := net.Dial("tcp", *forward)
			if err != nil {
				log.Println(err)
				c.Close()
				return
			}

			(&trafficConn{Conn: c, traffic: trafficLogger()}).Join(to)
		}(c)
	}
}
