package libtower

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

const DefaultHTTPSPort = "443"

// HTTPS type
type HTTPS struct {
	Host    string
	Port    string
	Timeout time.Duration

	Start    time.Time
	End      time.Time
	Duration time.Duration
}

// HTTPSCheck checks tls certificate is valid
func (hs *HTTPS) HTTPSCheck(ctx context.Context) (bool, time.Time, error) {
	if hs.Port == "" {
		hs.Port = DefaultHTTPSPort
	}
	address := hs.Host + ":" + hs.Port
	dialer := net.Dialer{
		Timeout: hs.Timeout,
	}
	hs.Start = time.Now()
	conn, err := tls.DialWithDialer(&dialer, "tcp", address, &tls.Config{
		InsecureSkipVerify: false,
	})
	hs.End = time.Now()
	hs.Duration = hs.End.Sub(hs.Start)
	if err != nil {
		return false, time.Time{}, err
	}
	if conn != nil {
		defer conn.Close()
		var NotAfter = conn.ConnectionState().PeerCertificates[0].NotAfter
		for _, cert := range conn.ConnectionState().PeerCertificates {
			if cert.NotAfter.Before(NotAfter) {
				NotAfter = cert.NotAfter
			}
		}
		return true, NotAfter, nil
	}
	return false, time.Time{}, nil
}
