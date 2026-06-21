package net

import (
	"solod.dev/so/c"
	"solod.dev/so/mem"
	"solod.dev/so/net/netip"
)

// sockAddr returns a *sockaddr view of the address storage.
func (stor *sockaddr_storage) sockAddr() *sockaddr {
	return c.PtrAs[sockaddr](stor)
}

// addrPort decodes the sockaddr in stor into a netip.AddrPort.
// If the family is not recognized, returns the zero AddrPort
// (whose Addr is invalid).
func (stor *sockaddr_storage) addrPort() netip.AddrPort {
	base := stor.sockAddr()
	if base.sa_family == c_AF_INET {
		s4 := c.PtrAs[sockaddr_in](stor)
		var ip [4]byte
		mem.Copy(&ip[0], &s4.sin_addr, 4)
		ipAddr := netip.AddrFromSlice(ip[:])
		port := ntohs(s4.sin_port)
		return netip.AddrPortFrom(ipAddr, port)
	}
	if base.sa_family == c_AF_INET6 {
		s6 := c.PtrAs[sockaddr_in6](stor)
		var ip [16]byte
		mem.Copy(&ip[0], &s6.sin6_addr, 16)
		ipAddr := netip.AddrFromSlice(ip[:])
		port := ntohs(s6.sin6_port)
		return netip.AddrPortFrom(ipAddr, port)
	}
	return netip.AddrPort{}
}

// fill encodes ap into stor as a sockaddr_in or sockaddr_in6 and returns
// its length. If ap's IP is invalid (neither IPv4 nor IPv6), fill does
// nothing and returns 0.
func (stor *sockaddr_storage) fill(ap netip.AddrPort) c.UInt {
	var ipbuf [16]byte
	ip := ap.Addr().AsSlice(ipbuf[:])
	port := ap.Port()
	mem.Clear(stor, c.Sizeof[sockaddr_storage]())
	if len(ip) == 4 {
		s4 := c.PtrAs[sockaddr_in](stor)
		s4.sin_family = c_AF_INET
		s4.sin_port = htons(port)
		mem.Copy(&s4.sin_addr, &ip[0], 4)
		return c.UInt(c.Sizeof[sockaddr_in]())
	}
	if len(ip) == 16 {
		s6 := c.PtrAs[sockaddr_in6](stor)
		s6.sin6_family = c_AF_INET6
		s6.sin6_port = htons(port)
		mem.Copy(&s6.sin6_addr, &ip[0], 16)
		return c.UInt(c.Sizeof[sockaddr_in6]())
	}
	return 0
}

// maxUnixPath is the maximum length of a Unix socket path,
// including the NUL terminator.
const maxUnixPath = 104

// fillUnix encodes name into stor as a sockaddr_un and returns its length, or 0
// if name is empty or too long to fit (with its NUL) in sun_path. The length is
// always the full sockaddr_un size, which is correct for pathname sockets.
func (stor *sockaddr_storage) fillUnix(name string) c.UInt {
	if len(name) == 0 || len(name) >= maxUnixPath {
		return 0
	}
	mem.Clear(stor, c.Sizeof[sockaddr_storage]())
	sun := c.PtrAs[sockaddr_un](stor)
	sun.sun_family = c_AF_UNIX
	mem.Copy(&sun.sun_path[0], c.CString(name), len(name)) // NUL already from Clear
	return c.UInt(c.Sizeof[sockaddr_un]())
}

// unixName decodes the NUL-terminated path from a sockaddr_un in stor, as
// filled in by recvfrom for a datagram's source. An unbound (anonymous) peer
// has an empty path, which yields "". The returned string aliases stor.
func (stor *sockaddr_storage) unixName() string {
	sun := c.PtrAs[sockaddr_un](stor)
	return c.String(c.PtrAs[c.Char](&sun.sun_path[0]))
}
