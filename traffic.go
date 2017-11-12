package main

import (
	"io"
	"net"
	"sync"
)

type traffic interface {
	TraffIn(b []byte)
	TraffOut(b []byte)
}

type trafficConn struct {
	traffic
	net.Conn
}

func (tc *trafficConn) Read(b []byte) (n int, err error) {
	n, err = tc.Conn.Read(b)
	if n > 0 {
		tc.TraffIn(b[:n])
	}
	return
}

func (tc *trafficConn) Write(b []byte) (n int, err error) {
	n, err = tc.Conn.Write(b)
	if n > 0 {
		tc.TraffOut(b[:n])
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
