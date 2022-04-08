package main

import (
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
	proxy  *httputil.ReverseProxy
	target *url.URL
}

type CookieValue struct {
	Value string `json:"value"`
	Place string `json:"place"`
}

type ResponseMessage struct {
	Status  int                     `json:"status"`
	Message string                  `json:"message"`
	Data    *map[string]CookieValue `json:"data,omitempty"`
}
