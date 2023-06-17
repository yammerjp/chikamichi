package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseJSON struct {
}

func main() {
	http.HandleFunc("/v1/completion", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "This-is-OpenAI-API-Key" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{
    "error": {
        "message": "You didn't provide an API key. You need to provide your API key in an Authorization header using Bearer auth (i.e. Authorization: Bearer YOUR_KEY), or as the password field (with blank username) if you're accessing the API from your browser and are prompted for a username and password. You can obtain an API key from https://platform.openai.com/account/api-keys.",
        "type": "invalid_request_error",
        "param": null,
        "code": null
    }
}`))
			return
		}
		body := make([]byte, r.ContentLength)
		length, err := r.Body.Read(body)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//parse json
		var jsonBody map[string]interface{}
		err = json.Unmarshal(body[:length], &jsonBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("%v\n", jsonBody)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "choices": [{
    "index": 0,
    "message": {
      "role": "assistant",
      "content": "\n\nHello there, how may I assist you today?",
    },
    "finish_reason": "stop"
  }],
  "usage": {
    "prompt_tokens": 9,
    "completion_tokens": 12,
    "total_tokens": 21
  }
}`))
	})

	fmt.Println("example web server is listen :8001")
	http.ListenAndServe(":8001", nil)
}
