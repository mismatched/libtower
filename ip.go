package libtower

import (
	"net"
)

func isIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

func isIPv6(ip net.IP) bool {
	return !isIPv4(ip) && len(ip) == net.IPv6len
}
