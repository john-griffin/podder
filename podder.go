package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Starting on port %s ....", port)
	http.ListenAndServe(port, proxy())
}

func proxy() *httputil.ReverseProxy {
	proxyUrl, _ := url.Parse(os.Getenv("PROXY_URL"))
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		req.Header.Set("X-Proxy-Host", req.Host)
		req.Host = proxyUrl.Host
		req.SetBasicAuth(os.Getenv("USER"), os.Getenv("PASS"))
		director(req)
	}
	proxy.Transport = &proxyTransport{}
	return proxy
}

type proxyTransport struct{}

func (t *proxyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)

	if response.Header.Get("Content-Type") != "audio/mpeg" {
		body := new(bytes.Buffer)
		body.ReadFrom(response.Body)

		target := request.Host
		replacement := request.Header.Get("X-Proxy-Host")
		proto := request.Header.Get("X-Forwarded-Proto")

		if len(proto) > 0 {
			target = "http://" + target
			replacement = proto + "://" + replacement
		}

		bod := strings.Replace(body.String(), target, replacement, -1)
		buf := bytes.NewBufferString(bod)
		contentLength := strconv.Itoa(buf.Len())

		response.Body = ioutil.NopCloser(buf)
		response.Header.Set("Content-Length", contentLength)
	}

	return response, err
}
