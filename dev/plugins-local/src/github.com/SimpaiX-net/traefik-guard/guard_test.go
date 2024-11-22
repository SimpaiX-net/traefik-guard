package traefik_guard_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	guard "github.com/SimpaiX-net/traefik-guard"
)

func TestGuard(t *testing.T) {
	config := guard.CreateConfig()

	config.IPHeaders = []string{"cf-connecting-ip", "X-Forwarded-For"}
	config.TTL = time.Duration(time.Hour * 6).String()
	config.Timeout = time.Duration(time.Millisecond * 200).String()

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	
	guard, err := guard.New(context.Background(), next, config, "guard@test")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("cf-connecting-ip", "1.1.1.1")

	recorder := httptest.NewRecorder()
	guard.ServeHTTP(recorder, req)

	for k, v := range req.Header {
		if len(v) == 0 { continue }
		t.Logf("%s: %v", k, v[0])
	}
}

func BenchmarkGuard(b *testing.B) {
	config := guard.CreateConfig()

	config.IPHeaders = []string{"cf-connecting-ip", "X-Forwarded-For"}
	config.TTL = time.Duration(time.Hour * 6).String()
	config.Timeout = time.Duration(time.Millisecond * 200).String()

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
	
	guard, err := guard.New(context.Background(), next, config, "guard@test")
	if err != nil {
		b.Fatal(err)
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
			if err != nil {
				b.Logf("error: %s", err)
				continue 
			}

			req.Header.Add("cf-connecting-ip", "1.1.1.1")

			recorder := httptest.NewRecorder()
			guard.ServeHTTP(recorder, req)
		}
	})
}