package libtower

import (
	"context"
	"testing"
	"time"
)

func TestHTTPS_HTTPSCheck(t *testing.T) {
	ctx := context.TODO()

	hs := HTTPS{Host: "google.com", Port: "443", Timeout: time.Second * 2}
	_, expire, err := hs.HTTPSCheck(ctx)
	if err != nil {
		t.Errorf("test failed %v", err)
	}
	t.Logf("current ssl cert of %s will expire  in %s", hs.Host, expire.Format(time.RFC3339))

}
