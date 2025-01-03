############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache 'git=~2'

# Install dependencies
ENV GO111MODULE=on
WORKDIR $GOPATH/src/packages/goginapp/
COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/main .

############################
# STEP 2 build a small image
############################
FROM alpine:3

WORKDIR /

# Copy our static executable.
COPY --from=builder /go/main /go/main

# env on .env
# server pass to ssh
ARG SSHPASS 
# server username ssh
ARG SSH_USER
# server ip
ARG SSH_HOST
# port gin proxy-app
ARG PROXY_PORT_APP
# redirect port on server
ARG PROXY_PORT_SERVER
# lmm api port for redirect (gin)
ARG LLM_PORT_APP

ENV SSHPASS=${SSHPASS}
ENV SSH_USER=${SSH_USER}
ENV SSH_HOST=${SSH_HOST}
ENV PORT = ${PROXY_PORT_APP}
ENV PROXY_PORT_SERVER = ${PROXY_PORT_SERVER}
ENV LLM_PORT_APP=${LLM_PORT_APP}
ENV GIN_MODE release
EXPOSE ${PROXY_PORT_APP}

WORKDIR /go

# Run the Go Gin binary.
ENTRYPOINT ["/go/main"]

# scripts 
RUN mkdir -p /scripts
COPY open-ssh-tunnel.sh /scripts
WORKDIR /scripts
RUN chmod +x open-ssh-tunnel.sh
RUN ./open-ssh-tunnel.sh
