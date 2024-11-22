package traefik_guard

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SimpaiX-net/ipqs"
	"github.com/SimpaiX-net/traefik-guard/helpers"
)

const defaultUA = "Mozilla/5.0 (Linux; Android 13; SM-S901U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36"

// Config parameters
//
// Please consult README.md for documentation
type Config struct {
	Proxy     string   `json:"proxy,omitempty"`
	TTL       string   `json:"ttl,omitempty"`
	Timeout   string   `json:"timeout,omitempty"`
	IPHeaders []string `json:"ipHeaders,omitempty"`
}

// Creates a configuration
func CreateConfig() *Config {
	return &Config{}
}

type Guard struct {
	name      string
	timeout   time.Duration
	ipHeaders []string
	next      http.Handler
	ctx       context.Context
	client    *ipqs.Client
}

const (
	success = "success"
	danger  = "danger"
	unknown = "unknown"
)

// Creates new Guard instance
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	ttl, timeout, err := helpers.ComposeDurations(config.TTL, config.Timeout)
	if err != nil {
		return nil, err
	}

	client := ipqs.New().SetProxy(config.Proxy)

	err = client.Provision()
	if err != nil {
		return nil, err
	}

	return &Guard{
		name:      name,
		next:      next,
		client:    client,
		ctx:       context.WithValue(ctx, ipqs.TTL_key, ttl),
		timeout:   timeout,
		ipHeaders: config.IPHeaders,
	}, nil
}

func (g *Guard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer g.next.ServeHTTP(w, r)

	var lookupInHeader string
	for _, header := range g.ipHeaders {
		h := r.Header.Get(header)
		if h != "" {
			lookupInHeader = h
			break
		}
	}

	if lookupInHeader == "" {
		r.Header.Set("X-Guard-Success", "-1")
		r.Header.Set("X-Guard-Info", "IP header not found")
		return
	}

	os.Stdout.WriteString(fmt.Sprintf("[GUARD-SCAN-START]: {\"ip\": \"%s\"}\n", lookupInHeader))
	defer os.Stdout.WriteString(fmt.Sprintf("[GUARD-SCAN-END]: {\"ip\": \"%s\"}\n", lookupInHeader))

	ctx, cancel := context.WithTimeout(g.ctx, g.timeout)
	defer cancel()

	err := g.client.GetIPQS(ctx, lookupInHeader, defaultUA)

	r.Header.Set("X-Guard-Success", "1")
	switch err {
		case nil:
			r.Header.Set("X-Guard-Rate", "LEGIT")
		case ipqs.ErrBadIPRep:
			r.Header.Set("X-Guard-Rate", "DANGER")
		default:
			r.Header.Set("X-Guard-Success", "-1")
			r.Header.Set("X-Guard-Rate", "UNKNOWN")
	}
}

