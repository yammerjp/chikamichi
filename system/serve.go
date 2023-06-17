package system

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/golang-jwt/jwt/v5"
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

func Serve() {
	if len(os.Args) >= 2 && os.Args[1] == "--jwt-gen" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "my-auth-server",
			"sub": "john",
			"aud": "john",
		})
		s, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			panic(err)
		}
		fmt.Println(s)

		return
	}

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
