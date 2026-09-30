package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxy(target string) http.Handler {
	u, _ := url.Parse(target)
	p := httputil.NewSingleHostReverseProxy(u)
	d := p.Director
	p.Director = func(r *http.Request) { d(r); r.Host = u.Host }
	return p
}

func main() {
	http.Handle("/tvarkarastis/", proxy("https://keltas.lt"))
	http.Handle("/lt/tvarkarastis.php", proxy("http://www.kopos.lt"))
	http.Handle("/v1/", proxy("https://api.meteo.lt"))
	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
