// Bad code warning, I wrote this with no prior thought.
package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
)

var (
	region  = os.Getenv("FLY_REGION")
	tcpAddr = "0.0.0.0:4000"
	udpAddr = "fly-global-services:4001"
)

func buildReply(raddr string, rlen int) string {
	return fmt.Sprintf("region:%s,raddr:%s,len=%d\n", region, raddr, rlen)
}

func listenTCP() error {
	l, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		slog.Error("Unable to start the TCP listener", slog.String("err", err.Error()))
		return err
	}

	defer l.Close()
	buf := make([]byte, 9600)

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("Unable to accept the TCP connection", slog.String("err", err.Error()))
			continue
		}

		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("Unable to read from the TCP connection", slog.String("err", err.Error()))
			continue
		}

		_, err = conn.Write([]byte(buildReply(conn.RemoteAddr().String(), n)))
		if err != nil {
			slog.Error("Unable to write to the TCP connection", slog.String("err", err.Error()))
			continue
		}

		err = conn.Close()
		if err != nil {
			slog.Error("Unable to close the TCP connection", slog.String("err", err.Error()))
			continue
		}
	}
}

func listenUDP() error {
	l, err := net.ListenPacket("udp", udpAddr)
	if err != nil {
		slog.Error("Unable to start the UDP listener", slog.String("err", err.Error()))
		return err
	}

	defer l.Close()
	buf := make([]byte, 1500) // Won't use all 1500, but let's be conservative

	for {
		n, raddr, err := l.ReadFrom(buf)
		if err != nil {
			slog.Error("Unable to read from the UDP listener", slog.String("err", err.Error()))
			continue
		}

		_, err = l.WriteTo([]byte(buildReply(raddr.String(), n)), raddr)
		if err != nil {
			slog.Error("Unable to write to the UDP connection", slog.String("err", err.Error()))
			continue
		}
	}
}

func init() {
	// Yeah don't ever do this in "production." This code is supposed to be a
	// playground for me.
	if region == "" {
		udpAddr = "0.0.0.0:4001"
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})))

	go func() {
		err := listenUDP()
		if err != nil {
			panic(err)
		}
	}()

	err := listenTCP()
	if err != nil {
		panic(err)
	}
}
