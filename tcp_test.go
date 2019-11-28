package libtower

import (
	"context"
	"net/url"
	"testing"
	"time"
)

func TestTCPPortCheck(t *testing.T) {
	ctx := context.TODO()
	URL, err := url.Parse("tcp://google.com:80")
	if err != nil {
		t.Errorf("test failed %v", err)
	}
	tr := TCP{URL: URL, Timeout: time.Second * 2}
	_, err = tr.TCPPortCheck(ctx)
	if err != nil {
		t.Errorf("test failed %v", err)
	}
}
