package storage

import (
	"net/url"
	"strings"
	"syscall/js"
)

func Get(key string) (v string) {
	defer func() {
		if recover() != nil {
			v = getCookie(key)
		}
	}()
	if r := js.Global().Get("localStorage").Call("getItem", key); r.Type() == js.TypeString {
		return r.String()
	}
	return ""
}

func Set(key, value string) {
	defer func() {
		if recover() != nil {
			setCookie(key, value)
		}
	}()
	js.Global().Get("localStorage").Call("setItem", key, value)
}

func getCookie(key string) string {
	for _, c := range strings.Split(js.Global().Get("document").Get("cookie").String(), ";") {
		if k, v, ok := strings.Cut(strings.TrimSpace(c), "="); ok && k == key {
			v, _ = url.QueryUnescape(v)
			return v
		}
	}
	return ""
}

func setCookie(key, value string) {
	js.Global().Get("document").Set("cookie", key+"="+url.QueryEscape(value)+"; max-age=31536000; path=/; SameSite=Lax")
}
