VERSION := "0.0.1"

build:
	docker build --platform=linux/amd64 -t ghcr.io/colearendt/openai-proxy:{{VERSION}} .

push:
	docker push ghcr.io/colearendt/openai-proxy:{{VERSION}}

install:
	go install openai-proxy

run:
	go run main.go

deploy:
    kubectl apply -f deploy/openai-proxy.yaml

namespace:
    #!/bin/bash
    kubectl create namespace ai

secret:
    #!/bin/bash
    if [ -z "${OPENAI_API_KEY}" ]; then echo "OPENAI_API_KEY" not found!; exit 1; fi
    kubectl create secret -n ai generic openai-credential --from-literal=OPENAI_API_KEY=${OPENAI_API_KEY}
