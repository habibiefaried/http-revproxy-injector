package main

import (
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
	proxy  *httputil.ReverseProxy
	target *url.URL
}

type HTTPReqData struct {
	Value string `json:"value"`
	Place string `json:"place"`
}

type RequestMessage struct {
	Name  string `json:"name"`
	Value string `json:"value",omitempty`
	Place string `json:"place",omitempty`
}

type ResponseMessage struct {
	Message string                  `json:"message"`
	Data    *map[string]HTTPReqData `json:"data,omitempty"`
}
