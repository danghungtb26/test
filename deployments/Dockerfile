FROM golang:1.21-alpine as builder

LABEL stage=builder

ARG CGO_ENABLED=0
ARG GO111MODULE=on
#ARG GOPRIVATE=git.ctisoftware.vn
ARG GIT_USERNAME=
ARG GIT_ACCESS_TOKEN=

ARG service_name=

WORKDIR /go/src/app

RUN apk add --no-cache openssh-client git && \
    echo -e "[url \"https://${GIT_USERNAME}:${GIT_ACCESS_TOKEN}@git.ctisoftware.vn/\"]\n\tinsteadOf = git@git.ctisoftware.vn:" >> ~/.gitconfig && \
    echo -e "[url \"https://${GIT_USERNAME}:${GIT_ACCESS_TOKEN}@git.ctisoftware.vn/\"]\n\tinsteadOf = https://git.ctisoftware.vn/" >> ~/.gitconfig

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN cd ./cmd && \
    go build -ldflags "-s -w" -trimpath -o gift *.go

# SERVICE BUILDER
FROM alpine:latest

WORKDIR /opt/gift-service

ENV ADDR_GRAPH 0.0.0.0:8081
ENV ADDR_GRPC 0.0.0.0:9091

COPY --from=builder /go/src/app/cmd/gift .
COPY ./configs/.env .

RUN apk update && apk add libc6-compat

RUN chmod +x ./gift

CMD ./gift --configSource=.env serve --addr-graph=$ADDR_GRAPH --addr-grpc=$ADDR_GRPC

