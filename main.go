package plugin_request_x_url

import (
	"context"
	"net/http"
	"strings"
)

const defaultHeader = "X-Request-Url"

type Config struct {
	HeaderName string
}

func CreateConfig() *Config {
	return &Config{
		HeaderName: defaultHeader,
	}
}

func New(ctx context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// --- scheme ---
		scheme := firstHeaderValue(r.Header, "X-Forwarded-Proto")
		if scheme == "" {
			if r.TLS != nil {
				scheme = "https"
			} else {
				scheme = "http"
			}
		}

		// --- host ---
		host := firstHeaderValue(r.Header, "X-Forwarded-Host")
		if host == "" {
			host = r.Host
		}

		// --- path + query ---
		uri := r.URL.RequestURI() // includes query string

		fullURL := scheme + "://" + host + uri

		// set (not add) чтобы не накапливались дубли при повторных проксях
		r.Header.Set(config.HeaderName, fullURL)

		next.ServeHTTP(rw, r)
	}), nil
}

func firstHeaderValue(h http.Header, key string) string {
	v := strings.TrimSpace(h.Get(key))
	if v == "" {
		return ""
	}
	// на всякий случай режем "https, http" или "a,b"
	if i := strings.IndexByte(v, ','); i >= 0 {
		v = v[:i]
	}
	return strings.TrimSpace(v)
}
