# syntax=docker/dockerfile:1
FROM golang:1.22.0-alpine3.19 as base



# ---------- ENVS
ENV DOCKER_CONTENT_TRUST=1



# ---------- BUILD
FROM base as build
WORKDIR /app
COPY . /app/
RUN apk add --no-cache ca-certificates && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-api /app/cmd/app



# ---------- MAIN
# FROM scratch
WORKDIR /app
# COPY --from=build /app/weather-api .
ENTRYPOINT [ "./weather-api" ]
