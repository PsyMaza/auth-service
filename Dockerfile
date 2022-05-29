# Builder

ARG GITHUB_PATH=gitlab.com/g6834/team17/auth_service

FROM golang:1.18-alpine AS builder

WORKDIR /home/${GITHUB_PATH}

RUN apk add --update make git protoc protobuf protobuf-dev curl
COPY Makefile Makefile
RUN make deps-go
COPY . .
RUN make build

# gRPC Server

FROM alpine:latest as server
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/auth-server .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .
#COPY --from=builder /home/${GITHUB_PATH}/migrations/ ./migrations

RUN chown root:root auth-server

EXPOSE 50051
EXPOSE 8080
EXPOSE 8083

CMD ["./auth-server"]
