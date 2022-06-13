# Builder
ARG GITHUB_PATH=gitlab.com/g6834/team17/auth-service

FROM golang:1.18-alpine as builder

WORKDIR /app/${GITHUB_PATH}

RUN apk add --update make git curl
COPY Makefile Makefile
COPY . .
RUN make build

# Auth server
FROM alpine:latest as server
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
WORKDIR /root/

COPY --from=builder /app/${GITHUB_PATH}/bin/auth-service .
COPY --from=builder /app/${GITHUB_PATH}/config.yaml .
COPY --from=builder /app/${GITHUB_PATH}/users.yaml .
COPY --from=builder /app/${GITHUB_PATH}/migrations/ ./migrations

RUN chown root:root app-services

EXPOSE 3000

CMD ["./auth-service"]
