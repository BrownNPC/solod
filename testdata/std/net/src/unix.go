package main

import (
	"solod.dev/so/io"
	"solod.dev/so/net"
	"solod.dev/so/os"
	"solod.dev/so/time"
)

// unixDir is a temporary directory holding the test socket files.
var unixDir string

// unixDirBuf backs unixDir; the path is a view into it.
var unixDirBuf [256]byte

func testUnix() {
	dir, err := os.MkdirTemp(unixDirBuf[:], "", "so-net-unix")
	noError(err)
	unixDir = dir

	testUnix_Resolve()
	testUnix_StreamDial()
	testUnix_StreamReadEOF()
	testUnix_DialRefused()
	testUnix_Datagram()
	testUnix_ReadDeadline()
	testUnix_CloseErrors()
	testUnix_UnlinkOnClose()

	noError(os.Remove(unixDir))
}

func testUnix_Resolve() {
	print("- Unix resolve...")
	addr, err := net.ResolveUnixAddr("unix", "/tmp/echo.sock")
	noError(err)
	if addr.Name != "/tmp/echo.sock" || addr.Net != "unix" {
		panic("unexpected ResolveUnixAddr result")
	}
	if addr.Network() != "unix" || addr.String() != "/tmp/echo.sock" {
		panic("unexpected UnixAddr Network/String")
	}

	gram, err := net.ResolveUnixAddr("unixgram", "/tmp/dg.sock")
	noError(err)
	if gram.Net != "unixgram" || gram.Network() != "unixgram" {
		panic("unexpected unixgram network")
	}

	// unixpacket is intentionally unsupported, as is any other network.
	if _, err := net.ResolveUnixAddr("unixpacket", "/tmp/x.sock"); err != net.ErrUnknownNetwork {
		panic("unixpacket should be unknown")
	}
	if _, err := net.ResolveUnixAddr("bogus", "/tmp/x.sock"); err != net.ErrUnknownNetwork {
		panic("bogus network should be unknown")
	}
	println("ok")
}

func testUnix_StreamDial() {
	print("- Unix stream dial...")
	// A single-threaded loopback echo: the connect queues into the listener
	// backlog, so Accept does not block on another thread.
	var pathBuf [320]byte
	laddr, err := net.ResolveUnixAddr("unix", unixPath(pathBuf[:], "stream.sock"))
	noError(err)
	ln, err := net.ListenUnix("unix", &laddr)
	noError(err)
	if ln.Addr().Name != laddr.Name {
		panic("listener addr mismatch")
	}

	raddr := ln.Addr()
	client, err := net.DialUnix("unix", nil, &raddr)
	noError(err)
	if client.RemoteAddr().Name != raddr.Name {
		panic("client remote addr mismatch")
	}

	server, err := ln.Accept()
	noError(err)
	if server.LocalAddr().Name != raddr.Name {
		panic("accepted local addr mismatch")
	}

	// Client writes, server echoes, client reads it back.
	if _, err := client.Write([]byte("ping")); err != nil {
		panic(err)
	}
	var buf [256]byte
	n, err := server.Read(buf[:])
	noError(err)
	if _, err := server.Write(buf[:n]); err != nil {
		panic(err)
	}
	var got [256]byte
	n, err = client.Read(got[:])
	noError(err)
	if string(got[:n]) != "ping" {
		panic("echo mismatch")
	}

	client.Close()
	server.Close()
	noError(ln.Close())
	println("ok")
}

func testUnix_StreamReadEOF() {
	print("- Unix stream read EOF...")
	// Connect a pair, then close the server end; the client's next read must
	// report end of stream.
	var pathBuf [320]byte
	laddr, err := net.ResolveUnixAddr("unix", unixPath(pathBuf[:], "eof.sock"))
	noError(err)
	ln, err := net.ListenUnix("unix", &laddr)
	noError(err)
	raddr := ln.Addr()
	client, err := net.DialUnix("unix", nil, &raddr)
	noError(err)
	server, err := ln.Accept()
	noError(err)

	noError(server.Close())
	var buf [16]byte
	if _, err := client.Read(buf[:]); err != io.EOF {
		panic("expected EOF")
	}

	client.Close()
	noError(ln.Close())
	println("ok")
}

func testUnix_DialRefused() {
	print("- Unix dial refused...")
	// Dialing a path with no socket file (nothing listening) must fail.
	var pathBuf [320]byte
	addr, err := net.ResolveUnixAddr("unix", unixPath(pathBuf[:], "refused.sock"))
	noError(err)
	if _, err := net.DialUnix("unix", nil, &addr); err == nil {
		panic("expected dial to a missing socket to fail")
	}
	println("ok")
}

func testUnix_Datagram() {
	print("- Unix datagram...")
	// Two bound datagram sockets exchange messages in both directions, each
	// receiver checking the reported source path against the sender's address.
	var pathBufA [320]byte
	var pathBufB [320]byte
	addrA, err := net.ResolveUnixAddr("unixgram", unixPath(pathBufA[:], "dga.sock"))
	noError(err)
	a, err := net.ListenUnixgram("unixgram", &addrA)
	noError(err)

	addrB, err := net.ResolveUnixAddr("unixgram", unixPath(pathBufB[:], "dgb.sock"))
	noError(err)
	b, err := net.ListenUnixgram("unixgram", &addrB)
	noError(err)

	// A -> B.
	bAddr := b.LocalAddr()
	if _, err := a.WriteTo([]byte("ping"), &bAddr); err != nil {
		panic(err)
	}
	var buf [256]byte
	r, err := b.ReadFrom(buf[:])
	noError(err)
	if string(buf[:r.N]) != "ping" {
		panic("A->B payload mismatch")
	}
	if r.Addr.Name != a.LocalAddr().Name {
		panic("A->B source addr mismatch")
	}

	// B -> A, replying to the learned source address.
	if _, err := b.WriteTo([]byte("pong"), &r.Addr); err != nil {
		panic(err)
	}
	var buf2 [256]byte
	r2, err := a.ReadFrom(buf2[:])
	noError(err)
	if string(buf2[:r2.N]) != "pong" {
		panic("B->A payload mismatch")
	}
	if r2.Addr.Name != b.LocalAddr().Name {
		panic("B->A source addr mismatch")
	}

	noError(a.Close())
	noError(b.Close())
	println("ok")
}

func testUnix_ReadDeadline() {
	print("- Unix read deadline...")
	// A ReadFrom with a short deadline and no data must time out.
	var pathBuf [320]byte
	laddr, err := net.ResolveUnixAddr("unixgram", unixPath(pathBuf[:], "dl.sock"))
	noError(err)
	conn, err := net.ListenUnixgram("unixgram", &laddr)
	noError(err)

	noError(conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond)))
	var buf [16]byte
	if _, err := conn.ReadFrom(buf[:]); err != net.ErrTimeout {
		panic("expected timeout")
	}

	noError(conn.Close())
	println("ok")
}

func testUnix_CloseErrors() {
	print("- Unix close errors...")
	// A double close, and any I/O after close, must report ErrClosed.
	var pathBuf [320]byte
	laddr, err := net.ResolveUnixAddr("unixgram", unixPath(pathBuf[:], "close.sock"))
	noError(err)
	conn, err := net.ListenUnixgram("unixgram", &laddr)
	noError(err)

	noError(conn.Close())
	if err := conn.Close(); err != net.ErrClosed {
		panic("expected ErrClosed on double close")
	}
	var buf [16]byte
	if _, err := conn.ReadFrom(buf[:]); err != net.ErrClosed {
		panic("expected ErrClosed on ReadFrom after close")
	}
	if _, err := conn.WriteTo(buf[:], &laddr); err != net.ErrClosed {
		panic("expected ErrClosed on WriteTo after close")
	}
	println("ok")
}

func testUnix_UnlinkOnClose() {
	print("- Unix unlink on close...")
	// Listening creates the socket file; Close must remove it. After Close, the
	// path is gone, so removing it again reports "not exist".
	var pathBuf [320]byte
	laddr, err := net.ResolveUnixAddr("unix", unixPath(pathBuf[:], "unlink.sock"))
	noError(err)
	ln, err := net.ListenUnix("unix", &laddr)
	noError(err)
	noError(ln.Close())

	if err := os.Remove(laddr.Name); err != os.ErrNotExist {
		panic("socket file should have been unlinked on Close")
	}
	println("ok")
}

// unixPath builds unixDir + "/" + name into buf and returns it.
func unixPath(buf []byte, name string) string {
	b := buf[:0]
	b = append(b, unixDir...)
	b = append(b, '/')
	b = append(b, name...)
	return string(b)
}
