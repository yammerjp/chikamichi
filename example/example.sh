#!/bin/sh

set -ex

cd "$(dirname $0)"

cd destination
pwd
go run main.go &

cd ../issuer
pwd
go run main.go &

cd ../..
pwd
export PORT=8800
export OPENAI_API_KEY="This-is-OpenAI-API-Key"
go run main.go serve &
		// req.URL.Scheme = "https"
