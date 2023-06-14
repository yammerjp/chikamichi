package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

const proxyHost = "api.openai.com"

func getModifyRequest() func(req *http.Request) {
	return func(req *http.Request) {
		req.Header = http.Header{}
		req.URL.Scheme = "https"
		req.URL.Host = proxyHost
		req.Host = proxyHost
		req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
		req.Header.Set("Content-Type", "application/json")
	}
}

func getModifyResponse() func(*http.Response) error {
	return func(res *http.Response) error {
		res.Header = http.Header{}
		appendAccessControllAllowHeaders(&res.Header)
		return nil
	}
}

func main() {
	rp := &httputil.ReverseProxy{Director: getModifyRequest()}
	rp.ModifyResponse = getModifyResponse()
	server := http.Server{
		Addr:    ":8081",
		Handler: WithLogging(WithPreflight(WithAuthorization(rp))),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
