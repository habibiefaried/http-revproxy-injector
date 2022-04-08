package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	port         = flag.Int("port", 4321, "Running port")
	host         = flag.String("host", "", "Target host")
	injecteddata map[string]CookieValue
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

	injecteddata = map[string]CookieValue{}

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

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		m.RLock()
		defer func() {
			m.RUnlock()
		}()
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "OK", Data: &injecteddata})
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("%v", err)})
			return
		}
		var t RequestMessage
		err = json.Unmarshal(body, &t)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("%v", err)})
			return
		}

		name := t.Name
		value := t.Value
		place := t.Place

		if (name == "") || (value == "") {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(ResponseMessage{Message: "variable name and value must exist"})
			return
		}

		if place == "" {
			place = "header"
		} else {
			w.WriteHeader(402)
			if (place != "header") && (place != "form") && (place != "query") {
				json.NewEncoder(w).Encode(ResponseMessage{Message: "variable 'place' is incorrect, should be 'header' or 'form' or 'query'"})
				return
			}
		}

		m.Lock()
		injecteddata[name] = CookieValue{
			Value: value,
			Place: place,
		}
		m.Unlock()

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Data is injected"})
	default:
		w.WriteHeader(504)
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-ProxyInjector", "In Action")
	r.URL.Scheme = ph.target.Scheme
	r.URL.Host = ph.target.Host
	r.Host = ph.target.Host
	ph.proxy.ServeHTTP(w, r)
}
