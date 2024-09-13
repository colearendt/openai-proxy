# from https://stackoverflow.com/a/46532352/6570011
# build stage
FROM golang:1.21
LABEL org.opencontainers.image.source=https://github.com/colearendt/openai-proxy

ADD . /src
RUN set -x && \
    cd /src && \
    CGO_ENABLED=0 GOOS=linux go build -a -o openai-proxy

# final stage
FROM alpine
WORKDIR /app
COPY --from=0 /src/openai-proxy /app/
ENTRYPOINT /app/openai-proxy
EXPOSE 8080
