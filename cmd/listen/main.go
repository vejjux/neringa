package main

import (
	"bufio"
	"errors"
	"io/fs"
	"log"
	"maps"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/netip"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	ttl      = 30 * time.Minute
	errorTTL = time.Minute
)

type entry struct {
	status int
	header http.Header
	body   []byte
	until  time.Time
}

func proxy(target string) http.Handler {
	u, _ := url.Parse(target)
	p := httputil.NewSingleHostReverseProxy(u)
	d := p.Director
	p.Director = func(r *http.Request) {
		d(r)
		r.Host = u.Host
		r.Header.Del("Accept-Encoding")
	}
	return cache(p)
}

func cache(h http.Handler) http.Handler {
	var mu, fetch sync.Mutex
	entries := make(map[string]entry)
	lookup := func(key string) (entry, bool) {
		mu.Lock()
		defer mu.Unlock()
		e, ok := entries[key]
		return e, ok && time.Now().Before(e.until)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			h.ServeHTTP(w, r)
			return
		}

		key := r.URL.String()
		e, ok := lookup(key)
		if !ok {
			fetch.Lock()
			if e, ok = lookup(key); !ok {
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, r)
				e = entry{status: rec.Code, header: rec.Header(), body: rec.Body.Bytes(), until: time.Now().Add(errorTTL)}
				if rec.Code == http.StatusOK {
					e.until = time.Now().Add(ttl)
				}
				mu.Lock()
				entries[key] = e
				mu.Unlock()
			}
			fetch.Unlock()
		}

		maps.Copy(w.Header(), e.header)
		w.WriteHeader(e.status)
		w.Write(e.body)
	})
}

// loadAllow reads one IP or CIDR per line; a missing file means allow all.
func loadAllow(path string) ([]netip.Prefix, error) {
	f, err := os.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var allow []netip.Prefix
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if a, err := netip.ParseAddr(line); err == nil {
			allow = append(allow, netip.PrefixFrom(a.Unmap(), a.Unmap().BitLen()))
			continue
		}
		p, err := netip.ParsePrefix(line)
		if err != nil {
			return nil, errors.New(path + ": invalid entry " + line)
		}
		allow = append(allow, p.Masked())
	}
	return allow, s.Err()
}

type allowListener struct {
	net.Listener
	allow []netip.Prefix
}

func (l allowListener) Accept() (net.Conn, error) {
	for {
		c, err := l.Listener.Accept()
		if err != nil {
			return nil, err
		}
		if ap, err := netip.ParseAddrPort(c.RemoteAddr().String()); err == nil && l.allowed(ap.Addr().Unmap()) {
			return c, nil
		}
		c.Close()
	}
}

func (l allowListener) allowed(a netip.Addr) bool {
	for _, p := range l.allow {
		if p.Contains(a) {
			return true
		}
	}
	return false
}

func main() {
	http.Handle("/tvarkarastis/", proxy("https://keltas.lt"))
	http.Handle("/lt/tvarkarastis.php", proxy("http://www.kopos.lt"))
	http.Handle("/v1/", proxy("https://api.meteo.lt"))
	http.Handle("/", http.FileServer(http.Dir("web")))
	addr := "localhost:8080"
	if a := os.Getenv("NERINGA_LISTEN_ADDR"); a != "" {
		addr = a
	}

	allow, err := loadAllow("allow.txt")
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	if len(allow) > 0 {
		log.Printf("allowing %d entries from allow.txt", len(allow))
		ln = allowListener{ln, allow}
	} else {
		log.Print("allowing all connections")
	}
	log.Fatal(http.ServeTLS(ln, nil, "cert.pem", "key.pem"))
}
