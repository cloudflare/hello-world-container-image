FROM golang:1.23 AS build

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download && go mod verify

COPY main.go index.html .
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/app ./...

FROM debian:bookworm-slim
COPY --from=build /usr/local/bin/app /usr/local/bin/

CMD ["/usr/local/bin/app", "-addr", "0.0.0.0"]
