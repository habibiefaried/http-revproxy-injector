package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	port         = flag.Int("port", 4321, "Running port")
	host         = flag.String("host", "", "Target host")
	injecteddata map[string]HTTPReqData
	m            sync.RWMutex
)

func main() {
	flag.Parse()
	if *port == 4321 {
		log.Println("\nDefault port used")
	}

	if *host == "" {
		flag.Usage()
		log.Fatal("\nHost cannot be empty")
	}

	injecteddata = map[string]HTTPReqData{}

	remote, err := url.Parse(*host)
	if err != nil {
		panic(err)
	}

	http.Handle("/", &ProxyHandler{proxy: httputil.NewSingleHostReverseProxy(remote), target: remote})
	http.HandleFunc("/revpr0xyconfig", ConfigHandler)

	fmt.Printf("Listening on port %v...\n", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	if err != nil {
		panic(err)
	}
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-ProxyInjector", "In Action")
	r.URL.Scheme = ph.target.Scheme
	r.URL.Host = ph.target.Host
	r.Host = ph.target.Host
	r.Header.Set("X-ProxyInjector", "In Action")
	m.RLock()
	for k, v := range injecteddata {
		if v.Place == "header" {
			r.Header.Set(k, v.Value)
		}
	}
	m.RUnlock()
	ph.proxy.ServeHTTP(w, r)
}
