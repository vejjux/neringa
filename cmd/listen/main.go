package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	u, _ := url.Parse("https://keltas.lt")
	p := httputil.NewSingleHostReverseProxy(u)
	d := p.Director
	p.Director = func(r *http.Request) { d(r); r.Host = u.Host }
	http.Handle("/tvarkarastis/", p)
	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
