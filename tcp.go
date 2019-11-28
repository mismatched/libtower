package libtower

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"time"
)

// TCP type
type TCP struct {
	URL     *url.URL
	Timeout time.Duration

	CertFile       string
	PrivateKeyFile string

	Start    time.Time
	End      time.Time
	Duration time.Duration
}

// TCPPortCheck checks if a tcp port is open
func (tr *TCP) TCPPortCheck(ctx context.Context) (bool, error) {
	tr.Start = time.Now()
	conn, err := net.DialTimeout("tcp", tr.URL.Host, tr.Timeout)
	tr.End = time.Now()
	tr.Duration = tr.End.Sub(tr.Start)
	if err != nil {
		return false, err
	}
	if conn != nil {
		defer conn.Close()
		return true, nil
	}
	return false, nil
}

// TLSPortCheck check if a scured tcp port is open
func (tr *TCP) TLSPortCheck(ctx context.Context) (bool, error) {
	cert, err := tls.LoadX509KeyPair(tr.CertFile, tr.PrivateKeyFile)
	if err != nil {
		return false, err
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	tr.Start = time.Now()
	conn, err := net.DialTimeout("tcp", tr.URL.Host, tr.Timeout)
	if err != nil {
		return false, err
	}

	tlsConn := tls.Client(conn, &config)
	err = tlsConn.Handshake()
	if err != nil {
		return false, err
	}
	tr.End = time.Now()
	tr.Duration = tr.End.Sub(tr.Start)

	return true, nil
}
