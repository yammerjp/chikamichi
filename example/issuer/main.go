package main

import (
	"fmt"
	"os"
	"encoding/json"
	"net/http"
	"github.com/golang-jwt/jwt/v5"
)

type ResponseJSON struct {
  Jwt string
}

func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
  })

  http.HandleFunc("/jwt", func (w http.ResponseWriter, r *http.Request) {
    jwt, err := genJwt()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    res, err := json.Marshal(ResponseJSON{jwt})

    w.Header().Set("Content-Type", "application/json")
    w.Write(res)
  })

  fmt.Println("example web server is listen :8001")
  http.ListenAndServe(":8001", nil)
}

func genJwt() (string , error) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "my-auth-server",
			"sub": "john",
			"aud": "john",
		})
		s, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
      return "", err
		}
    return s, nil
}
