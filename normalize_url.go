package main

import (
	"net/url"
	"path"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// normalize hostname - lowercase
	host := strings.ToLower(u.Hostname())

	// remove default ports
	port := u.Port()
	if (u.Scheme == "http" && port == "80") || (u.Scheme == "https" && port == "443") {
		port = ""
	}

	// remove www.
	host = strings.TrimPrefix(host, "www.")

	// reassemble port
	if port != "" {
		host = host + ":" + port
	}

	// normalize path
	p := strings.ToLower(path.Clean(u.Path))

	// remove trailing slash unless its root
	if p == "/" || p == "." {
		p = ""
	}

	p = strings.TrimSuffix(p, "/")

	// construct final string
	if p != "" {
		return host + p, nil
	}
	return host, nil
}
