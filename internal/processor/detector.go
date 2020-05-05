package processor

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// Detector is a type that detects peer information
type Detector struct {
	client *http.Client
}

// Host is a struct that represent peer data
type Host struct {
	Domain string
	Server string
	IPs    []string
}

// NewDetector returns a new instance of Detector
func NewDetector(client *http.Client) Detector {
	return Detector{
		client: client,
	}
}

// Detect returns the Peer data for the given hostname
func (d Detector) Detect(ctx context.Context, host string) (Host, error) {
	resp, err := d.tryHeadFirst(ctx, host)
	if err != nil {
		return Host{}, err
	}

	resp.Body.Close()

	var srv string

	h, ok := resp.Header["Server"]
	if ok && len(h) > 0 {
		srv = h[0]
		if idx := strings.IndexByte(srv, '/'); idx > 0 {
			srv = srv[:idx]
		}
	}

	return Host{
		Domain: host,
		Server: srv,
		IPs:    getIP(host),
	}, nil
}

func (d Detector) tryHeadFirst(ctx context.Context, host string) (*http.Response, error) {
	u := url.URL{
		Host:   host,
		Scheme: "https",
	}

	req, err := http.NewRequest("HEAD", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusMethodNotAllowed {
		req.Method = "GET"
		return d.client.Do(req)
	}

	return resp, nil
}

func getIP(host string) []string {
	res, err := net.LookupIP(host)
	if err != nil {
		return nil
	}

	ips := make([]string, len(res))

	for i, ip := range res {
		ips[i] = ip.String()
	}

	return ips
}
