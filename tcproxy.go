package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
)

var (
	addr    = flag.String("addr", ":9090", "tcproxy address")
	forward = flag.String("forward", "", "forwarding address")
)

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			to, err := net.Dial("tcp", *forward)
			if err != nil {
				log.Println(err)
				c.Close()
				return
			}
			wrapConn(c).Join(to)
		}(c)
	}
}

type trafficConn struct {
	net.Conn
	addr string
	id   int64
}

var _id int64

func wrapConn(c net.Conn) *trafficConn {
	return &trafficConn{
		Conn: c,
		addr: c.RemoteAddr().String(),
		id:   atomic.AddInt64(&_id, 1),
	}
}

func (tc *trafficConn) Read(b []byte) (n int, err error) {
	n, err = tc.Conn.Read(b)
	if n > 0 {
		log.Printf("[IN#%s#%d]%s\n", tc.addr, tc.id, string(b[:n]))
	}
	return
}

func (tc *trafficConn) Write(b []byte) (n int, err error) {
	n, err = tc.Conn.Write(b)
	if n > 0 {
		log.Printf("[OUT#%s#%d]%s\n", tc.addr, tc.id, string(b[:n]))
	}
	return
}

func (tc *trafficConn) Join(c net.Conn) {
	var wg sync.WaitGroup
	pipe := func(from, to net.Conn) {
		io.Copy(to, from)
		from.Close()
		to.Close()
		wg.Done()
	}
	wg.Add(2)
	go pipe(tc, c)
	go pipe(c, tc)
	wg.Wait()
}
