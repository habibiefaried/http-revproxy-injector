package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	port         = flag.Int("port", 4321, "Running port")
	host         = flag.String("host", "", "Target host")
	injecteddata map[string]CookieValue
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
		json.NewEncoder(w).Encode(ResponseMessage{Status: 200, Message: "OK", Data: &injecteddata})
	case "POST":
		name := r.FormValue("name")
		value := r.FormValue("value")
		place := r.FormValue("place")

		if (name == "") || (value == "") {
			json.NewEncoder(w).Encode(ResponseMessage{Status: 403, Message: "variable name and value must exist"})
			return
		}

		if place == "" {
			place = "cookie"
		} else {
			if (place != "cookie") && (place != "form") && (place != "query") {
				json.NewEncoder(w).Encode(ResponseMessage{Status: 403, Message: "variable 'place' is incorrect, should be 'cookie' or 'form' or 'query'"})
				return
			}
		}

		injecteddata[name] = CookieValue{
			Value: value,
			Place: place,
		}

		json.NewEncoder(w).Encode(ResponseMessage{Status: 201, Message: "Data is injected"})
	default:
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
