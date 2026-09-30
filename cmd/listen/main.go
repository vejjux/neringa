package main

import (
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
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

func main() {
	http.Handle("/tvarkarastis/", proxy("https://keltas.lt"))
	http.Handle("/lt/tvarkarastis.php", proxy("http://www.kopos.lt"))
	http.Handle("/v1/", proxy("https://api.meteo.lt"))
	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
