package libtower

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const (
	// ProtocolICMP DSCP
	IPv4ProtocolICMP = 1
	IPv6ProtocolICMP = 58
)

// Ping an address
// It needs root privileges to listen icmp on 0.0.0.0
func Ping(addr string, seq int) (*net.IPAddr, time.Duration, error) {
	// Resolve DNS and get the real IP of the it
	dst, _, err := DNSLookup(addr)
	if err != nil {
		return dst, 0, err
	}

	var c *icmp.PacketConn
	if isIPv4(dst.IP) {
		c, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			return nil, 0, err
		}
	} else if isIPv6(dst.IP) {
		c, err = icmp.ListenPacket("ip6:ipv6-icmp", "::")
		if err != nil {
			return nil, 0, err
		}
	} else {
		return nil, 0, fmt.Errorf("can not find version of ip(v4/v6):%s", dst.IP.String())
	}
	defer c.Close()

	var msg icmp.Message
	if isIPv4(dst.IP) {
		msg = icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  seq,
				Data: []byte(""),
			},
		}
	} else { //IPv6
		msg = icmp.Message{
			Type: ipv6.ICMPTypeEchoRequest,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  seq,
				Data: []byte(""),
			},
		}
	}
	bmsg, err := msg.Marshal(nil)
	if err != nil {
		return dst, 0, err
	}
	// Send ICMP message
	start := time.Now()
	n, err := c.WriteTo(bmsg, dst)
	if err != nil {
		return dst, 0, err
	} else if n != len(bmsg) {
		return dst, 0, fmt.Errorf("got %v; want %v", n, len(bmsg))
	}

	// Wait for an ICMP reply
	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return dst, 0, err
	}
	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return dst, 0, err
	}
	duration := time.Since(start)

	if isIPv4(dst.IP) {
		rm, err := icmp.ParseMessage(IPv4ProtocolICMP, reply[:n])
		if err != nil {
			return dst, 0, err
		}

		if rm.Type == ipv4.ICMPTypeEchoReply {
			return dst, duration, nil
		} else {
			return dst, 0, fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
		}
	} else { //IPv6
		rm, err := icmp.ParseMessage(IPv6ProtocolICMP, reply[:n])
		if err != nil {
			return dst, 0, err
		}
		if rm.Type == ipv6.ICMPTypeEchoReply {
			return dst, duration, nil
		} else {
			return dst, 0, fmt.Errorf("got %+v from %v; want echo reply v6", rm, peer)
		}
	}

}
